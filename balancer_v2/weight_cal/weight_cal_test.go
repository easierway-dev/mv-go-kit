package weight_cal

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	hslam_automic "github.com/hslam/atomic"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

func TestCalEWMA(t *testing.T) {
	Convey("TestCalEWMA", t, func() {
		//init adjuster
		adjuster := NewWeightAdjuster(0.9)
		//init counter
		now := time.Now().Unix()
		counter := &Counter{
			Timestamp:    now,
			SuccessCount: 1,
			TotalCount:   1,
			Vt:           hslam_automic.NewFloat64(1.0),
		}
		time.Sleep(time.Duration(1) * time.Second)
		now = time.Now().Unix()
		//CalEWMA
		adjuster.CalEWMA(now, counter)
		//check Vt
		So(counter.Vt.Load(), ShouldEqual, 1.0)
		//check SuccessCount
		So(counter.SuccessCount, ShouldEqual, 0)
	})
}

func TestGetWeightAdjusterAndNotify(t *testing.T) {
	Convey("TestGetWeightAdjusterAndNotify", t, func() {
		//init adjuster
		adjuster := NewWeightAdjuster(0.9)
		service1 := "1.1.1.1:1010"
		//notify succ
		adjuster.Notify(service1, balancer_common.Success)
		//sleep 1 sec
		time.Sleep(time.Duration(1) * time.Second)
		//notify failed
		adjuster.Notify(service1, balancer_common.Failed)
		//sleep 1 sec
		time.Sleep(time.Duration(1) * time.Second)
		adjuster.Notify(service1, balancer_common.Failed)
		//check weight
		weight := adjuster.GetWeight(service1)
		So(weight, ShouldEqual, 0.9)
	})
}

func TestClearEmptyCounter(t *testing.T) {
	Convey("TestClearEmptyCounter", t, func() {
		//init adjuster
		adjuster := NewWeightAdjuster(0.9)
		service1 := "1.1.1.1:1010"
		service2 := "2.2.2.2:9090"
		//notify
		adjuster.Notify(service2, balancer_common.Success)
		adjuster.Notify(service1, balancer_common.Success)
		//ClearEmptyCounter
		adjuster.ClearEmptyCounter(time.Duration(balancer_common.MaxTimeGap+1) * time.Second)
		for _, counter := range adjuster.counters {
			fmt.Println(*counter)
		}

		time.Sleep(time.Duration(balancer_common.MaxTimeGap+2) * time.Second)
		for _, counter := range adjuster.counters {
			fmt.Println(*counter)
		}
		//check
		So(len(adjuster.counters), ShouldEqual, 0)
	})
}
