package weight_cal

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

var stop bool

func RandomNotify(size int, key string, successRatio float64, interval time.Duration, adjuster *WeightAdjuster) {
	for i := 0; i < size; i++ {
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if rand.Float64() <= successRatio {
						adjuster.Notify(key, balancer_common.Success)
					} else {
						adjuster.Notify(key, balancer_common.Failed)
					}
				}
				if stop {
					return
				}
			}
		}()
	}
}

func RunAdjustGetWeight(adjuster *WeightAdjuster, service string) {
	for i := 0; i < 100; i++ {
		go func() {
			adjuster.GetWeight(service)
		}()
	}
}

func TestGetWeightAdjuster(t *testing.T) {
	Convey("GetWeightAdjusterTest", t, func() {
		adjuster := NewWeightAdjuster()

		service1 := "1.1.1.1:1010"
		service2 := "2.2.2.2:9090"
		RandomNotify(50, service1, 0.95, time.Duration(10)*time.Millisecond, adjuster)
		RandomNotify(60, service2, 0.90, time.Duration(10)*time.Millisecond, adjuster)

		i := 0
		fmt.Println("start TestGetWeightAdjuster")
		for {
			stop := false
			select {
			case <-time.Tick(time.Duration(2) * time.Second):
				fmt.Println("start range sec:", i)
				if i > 5 {
					stop = true
				}
				RunAdjustGetWeight(adjuster, service1)
				RunAdjustGetWeight(adjuster, service2)
				i += 1
				fmt.Println("second:", i, " service1_weight:", adjuster.GetWeight(service1))
				fmt.Println("second:", i, " service2_weight:", adjuster.GetWeight(service2))
			}
			if stop {
				break
			}
		}
		fmt.Println("end TestGetWeightAdjuster")
	})
}

func TestServiceWeightCul(t *testing.T) {
	Convey("TestServiceWeightCul", t, func() {
		adjuster := NewWeightAdjuster()

		service1 := "1.1.1.1:1010"
		service2 := "2.2.2.2:9090"
		RandomNotify(50, service1, 0.75, time.Duration(10)*time.Millisecond, adjuster)
		RandomNotify(60, service2, 0.60, time.Duration(10)*time.Millisecond, adjuster)
		stop = false

		i := 0
		fmt.Println("start TestServiceWeightCul")
		for {
			fmt.Println("start range sec:", i)
			select {
			case <-time.Tick(time.Duration(2) * time.Second):
				fmt.Println("start for ")
				if i > 20 {
					stop = true
				}
				i += 1
				ratio1 := adjuster.GetWeight(service1)
				now := time.Now().Unix()
				fmt.Println("end get weight 1 time:", now)
				ratio2 := adjuster.GetWeight(service2)
				now = time.Now().Unix()
				fmt.Println("end get weight 2 time:", now)
				fmt.Println("second:", i, "service1_real_ratio:", ratio1, " service1_cul_weight:", GetRatioByStep(ratio1, 0.05))
				fmt.Println("second:", i, "service2_real_ratio:", ratio2, " service2_cul_weight:", GetRatioByStep(ratio2, 0.05))
			}
			if stop {
				break
			}
		}
		fmt.Println("end TestServiceWeightCul")
	})
}

func TestZoneWeightCul(t *testing.T) {
	Convey("TestZoneWeightCul", t, func() {
		adjuster := NewWeightAdjuster()

		localZone := "us-aws-vg-1"
		otherZone := "us-aws-vg-2"
		stop = false
		RandomNotify(50, localZone, 0.87, time.Duration(10)*time.Millisecond, adjuster)
		RandomNotify(60, otherZone, 0.60, time.Duration(10)*time.Millisecond, adjuster)

		i := 0
		for {
			select {
			case <-time.Tick(time.Duration(2) * time.Second):
				fmt.Println("start range sec:", i)
				if i > 10 {
					stop = true
				}
				i += 1
				fmt.Println("second:", i, " localZone_cul_weight:", GetZoneWeight(adjuster, localZone, localZone))
				fmt.Println("second:", i, " otherZone_cul_weight:", GetZoneWeight(adjuster, localZone, otherZone))
			}
			if stop {
				break
			}
		}
	})
}
