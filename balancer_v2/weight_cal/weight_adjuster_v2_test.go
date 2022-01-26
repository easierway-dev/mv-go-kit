package weight_cal

import (
	"math/rand"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func RandomNotify(size int, key string, successRatio float64, interval time.Duration, adjuster *WeightAdjuster) {
	for i := 0; i < size; i++ {
		go func() {
			for range time.Tick(interval) {
				if rand.Float64() <= successRatio {
					adjuster.Notify(key, balancer_common.Success)
				} else {
					adjuster.Notify(key, balancer_common.Failed)
				}
			}
		}()
	}
}

func GetWeightAdjusterTest(t *testing.T) {
	adjuster := NewWeightAdjuster()

	RandomNotify(50, "1.1.1.1:1010", 0.95, time.Duration(10)*time.MilliceSecond, adjuster)
	RandomNotify(60, "2.2.2.2:9090", 0.90, time.Duration(10)*time.MilliceSecond, adjuster)

	Convey("GetWeightAdjusterTest", t, func() {
	})
}
