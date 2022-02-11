package weight_cal

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	//. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
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

func RunAdjustGetWeight(adjuster *WeightAdjuster, service string) {
	for i := 0; i < 100; i++ {
		go func() {
			adjuster.GetWeight(service)
		}()
	}
}

func TestGetWeightAdjuster(t *testing.T) {
	adjuster := NewWeightAdjuster()

	service1 := "1.1.1.1:1010"
	service2 := "2.2.2.2:9090"
	RandomNotify(50, service1, 0.95, time.Duration(10)*time.Millisecond, adjuster)
	RandomNotify(60, service2, 0.90, time.Duration(10)*time.Millisecond, adjuster)

	Convey("GetWeightAdjusterTest", t, func() {
		i := 0
		for range time.Tick(time.Duration(2) * time.Second) {
			if i > 5 {
				break
			}
			RunAdjustGetWeight(adjuster, service1)
			RunAdjustGetWeight(adjuster, service2)
			i += 1
			fmt.Println("second:", i, " service1_weight:", adjuster.GetWeight(service1))
			fmt.Println("second:", i, " service2_weight:", adjuster.GetWeight(service2))
		}
	})
}

func TestServiceWeightCul(t *testing.T) {
	adjuster := NewWeightAdjuster()

	service1 := "1.1.1.1:1010"
	service2 := "2.2.2.2:9090"
	RandomNotify(50, service1, 0.75, time.Duration(10)*time.Millisecond, adjuster)
	RandomNotify(60, service2, 0.60, time.Duration(10)*time.Millisecond, adjuster)

	Convey("TestServiceWeightCul", t, func() {
		i := 0
		for range time.Tick(time.Duration(2) * time.Second) {
			if i > 5 {
				break
			}
			i += 1
			ratio1 := adjuster.GetWeight(service1)
			ratio2 := adjuster.GetWeight(service2)
			fmt.Println("second:", i, "service1_real_ratio:", ratio1, " service1_cul_weight:", GetRatioByStep(ratio1, 0.05))
			fmt.Println("second:", i, "service2_real_ratio:", ratio2, " service2_cul_weight:", GetRatioByStep(ratio2, 0.05))
		}
	})
}

func TestZoneWeightCul(t *testing.T) {
	adjuster := NewWeightAdjuster()

	localZone := "us-aws-vg-1"
	otherZone := "us-aws-vg-2"
	RandomNotify(50, localZone, 0.87, time.Duration(10)*time.Millisecond, adjuster)
	RandomNotify(60, otherZone, 0.60, time.Duration(10)*time.Millisecond, adjuster)

	Convey("TestZoneWeightCul", t, func() {
		i := 0
		for range time.Tick(time.Duration(2) * time.Second) {
			if i > 10 {
				break
			}
			i += 1
			fmt.Println("second:", i, " localZone_cul_weight:", GetZoneWeight(adjuster, localZone, localZone))
			fmt.Println("second:", i, " otherZone_cul_weight:", GetZoneWeight(adjuster, localZone, otherZone))
		}
	})

}
