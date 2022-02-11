package weight_cal

import (
	"sync"
	"sync/atomic"
	"time"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

type WeightAdjuster struct {
	counters map[string]*Counter
	mutex    sync.RWMutex
}

type Counter struct {
	FailedCount  int32
	SuccessCount int32
	TotalCount   int32

	Vt        float64 //vt=βvt−1+(1−β)θt
	Timestamp int64
}

func NewWeightAdjuster() *WeightAdjuster {
	return &WeightAdjuster{
		counters: make(map[string]*Counter),
	}
}

func (adjuster *WeightAdjuster) Notify(key string, event int) {
	//init member
	now := time.Now().Unix()
	//check filed
	adjuster.mutex.RLock()
	counter, ok := adjuster.counters[key]
	adjuster.mutex.RUnlock()

	MaxTime := int64(60)
	if !ok || (now-counter.Timestamp > MaxTime) {
		counter = &Counter{
			Timestamp: now,
			Vt:        1.0, // init vt-1=1.0
		}
		adjuster.mutex.Lock()
		defer adjuster.mutex.Unlock()
		adjuster.counters[key] = counter
	}
	//EWMA:vt=βvt−1+(1−β)θt, β = 0.9
	beta := 0.9
	gama := 1 - beta
	timeGap := int(now - counter.Timestamp)
	if timeGap > 0 {
		Vt := 0.0
		/*for i := 0; i < timeGap-1; i++ {
			Vt = beta*counter.Vt + gama*(1)
		}*/
		totalCount := counter.TotalCount
		successCount := counter.SuccessCount
		if totalCount > 0 {
			successRate := float64(successCount) / float64(totalCount)
			if successRate >= 1.0 {
				successRate = 1.0
			}
			Vt = beta*counter.Vt + gama*(successRate)
		} else {
			Vt = beta*counter.Vt + gama*(1.0)
		}
		counter = &Counter{
			Timestamp: now,
			Vt:        Vt,
		}
		adjuster.mutex.Lock()
		defer adjuster.mutex.Unlock()
		adjuster.counters[key] = counter
	}
	//atomic add
	switch event {
	case balancer_common.Success:
		atomic.AddInt32(&counter.SuccessCount, 1)
	default:
		atomic.AddInt32(&counter.FailedCount, 1)
	}
	atomic.AddInt32(&counter.TotalCount, 1)
}

func (adjuster *WeightAdjuster) GetWeight(key string) float64 {
	adjuster.mutex.RLock()
	counter, ok := adjuster.counters[key]
	adjuster.mutex.RUnlock()
	if !ok {
		return 1
	}
	return counter.Vt
}
