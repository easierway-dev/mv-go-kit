package balancer_v2

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

type MyLogger struct {
}

func (logger *MyLogger) Infof(format string, v ...interface{}) {}
func (logger *MyLogger) Warnf(format string, v ...interface{}) {}

func RandomNotify(size int, serviceName string, zoneName string, successRatio float64, interval time.Duration, resolver *BalancerResolver) {
	for i := 0; i < size; i++ {
		go func() {
			for range time.Tick(interval) {
				if rand.Float64() <= successRatio {
					resolver.Notify(serviceName, zoneName, balancer_common.Success)
				} else {
					resolver.Notify(serviceName, zoneName, balancer_common.Failed)
				}
			}
		}()
	}
}

func Test_BalancerResolver(t *testing.T) {
	Convey("Test_BalancerResolver", t, func() {
		logger := &MyLogger{}
		//new resolver
		resolver, err := NewBalancerResolver(balancer_common.WeightedRoundRobin, balancer_common.TestingDiscover,
			"local_zone", "192.168.1.1:8500", "test_discover_service", time.Duration(2)*time.Second, logger)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		//new Notify
		RandomNotify(50, "192.168.1.1:10000", "local_zone", 0.99, time.Duration(10)*time.Millisecond, resolver)
		RandomNotify(50, "192.168.1.2:10000", "local_zone", 0.25, time.Duration(10)*time.Millisecond, resolver)
		RandomNotify(50, "192.168.1.3:10000", "local_zone", 0.25, time.Duration(10)*time.Millisecond, resolver)
		RandomNotify(50, "10.0.0.1:10000", "other_zone1", 0.99, time.Duration(10)*time.Millisecond, resolver)
		RandomNotify(50, "10.0.2.3:10000", "other_zone2", 0.98, time.Duration(10)*time.Millisecond, resolver)
		//discover node
		for j := 1; j <= 10; j++ {
			time.Sleep(time.Duration(3) * time.Second)
			countMap := make(map[string]int)
			for i := 0; i < 2000; i++ {
				node, err := resolver.DiscoverNode()
				if err == nil {
					countMap[node.Address] += 1
				} else {
					fmt.Println("discover err:", err)
				}
			}
			keys := []string{}
			for k, _ := range countMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				fmt.Println("second:", j, " ", key, ":", countMap[key])
			}
			fmt.Println("second end")
			fmt.Println("")
		}
	})
}
