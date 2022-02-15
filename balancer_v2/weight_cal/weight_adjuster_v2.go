package weight_cal

import (
	"sync"
	"sync/atomic"
	"time"

	hslam_automic "github.com/hslam/atomic"
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

	Vt        *hslam_automic.Float64 //vt=βvt−1+(1−β)θt
	Timestamp int64
}

func NewWeightAdjuster() *WeightAdjuster {
	return &WeightAdjuster{
		counters: make(map[string]*Counter),
	}
}

func (adjuster *WeightAdjuster) ClearEmptyCounter(interval time.Duration) {
	go func() {
		for range time.Tick(interval) {
			adjuster.mutex.Lock()
			defer adjuster.mutex.Unlock()
			for key, counter := range adjuster.counters {
				now := time.Now().Unix()
				if now-counter.Timestamp > balancer_common.MaxTimeGap {
					delete(adjuster.counters, key)
				}
			}
		}
	}()
}

func (adjuster *WeightAdjuster) CalEWMA(now int64, counter *Counter) {
	//EWMA:vt=βvt−1+(1−β)θt, β = 0.9
	timeGap := int(now - counter.Timestamp)
	if timeGap > 0 {
		beta := 0.9
		gama := 1 - beta
		Vt := 0.0
		totalCount := counter.TotalCount
		successCount := counter.SuccessCount
		if totalCount > 0 {
			successRate := float64(successCount) / float64(totalCount)
			if successRate >= 1.0 {
				successRate = 1.0
			}
			Vt = beta*counter.Vt.Load() + gama*(successRate)
		} else {
			Vt = beta*counter.Vt.Load() + gama*(1.0)
		}
		// use automic
		counter.Vt.CompareAndSwap(counter.Vt.Load(), Vt)
		atomic.CompareAndSwapInt64(&counter.Timestamp, counter.Timestamp, now)
		for {
			if atomic.CompareAndSwapInt32(&counter.SuccessCount, counter.SuccessCount, 0) {
				break
			}
		}
		for {
			if atomic.CompareAndSwapInt32(&counter.FailedCount, counter.FailedCount, 0) {
				break
			}
		}
		for {
			if atomic.CompareAndSwapInt32(&counter.TotalCount, counter.TotalCount, 0) {
				break
			}
		}
	}
}

func (adjuster *WeightAdjuster) Notify(key string, event int) {
	//init member
	now := time.Now().Unix()
	//check filed
	adjuster.mutex.RLock()
	counter, ok := adjuster.counters[key]
	adjuster.mutex.RUnlock()

	firstCreate := false
	if !ok || (now-counter.Timestamp > balancer_common.MaxTimeGap) {
		firstCreate = true
		counter = &Counter{
			Timestamp: now,
			Vt:        hslam_automic.NewFloat64(1.0), // init vt-1=1.0
		}
		adjuster.mutex.Lock()
		defer adjuster.mutex.Unlock()
		adjuster.counters[key] = counter
	}
	//cal EWMA
	if !firstCreate {
		adjuster.CalEWMA(now, counter)
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
	return counter.Vt.Load()
}
