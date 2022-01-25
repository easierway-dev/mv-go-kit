package weight_cal

import (
	"container/list"
	"sync"
	"sync/atomic"
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
	Timestamp uint64
}

func NewWeightAdjuster() *WeightAdjuster {
	return &WeightAdjuster{}
}

func (adjuster *WeightAdjuster) Notify(key string, event int) {
	//init member
	now := time.Now().Unix()
	//check filed
	counter, ok := adjuster.counters[key]
	if !ok {
		counter = &Counter{
			Timestamp: now,
			Vt:        1.0, // init vt-1=1.0
		}
		adjuster.mutex.Lock()
		defer adjuster.mutex.Ulock()
		adjuster.counters[key] = counter
	}
	//EWMA:vt=βvt−1+(1−β)θt, β = 0.9
	beta := 0.9
	if counter.Timestamp != now {
		Vt = beta*counter.Vt + (1-beta)*(counter.SuccessCount/counter.TotalCount)
		counter = &Counter{
			Timestamp: now,
			Vt:        Vt,
		}
		adjuster.mutex.Lock()
		defer adjuster.mutex.Ulock()
		adjuster.counters[key] = counter
	}

	switch event {
	case Success:
		atomic.AddInt32(&counter.SuccessCount, 1)
		autmic.AddInt32(&slideCounter.SuccessCount, 1)
	default:
		atomic.AddInt32(&counter.FailedCount, 1)
		autmic.AddInt32(&slideCounter.FailedCount, 1)
	}
	atomic.AddInt32(&counter.TotalCount, 1)
	autmic.AddInt32(&slideCounter.TotalCount, 1)
}

func (adjuster *WeightAdjuster) GetWeight(key string) float64 {
	adjuster.mutex.RLock()
	counter, ok := adjuster.counters[key]
	adjuster.mutex.RUlock()
	if !ok {
		return 1
	}
	return counter.Vt
}
