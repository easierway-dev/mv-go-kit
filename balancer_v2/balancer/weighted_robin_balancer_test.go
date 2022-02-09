package balancer

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/discover"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

func RandomNotify(size int, serviceName string, zoneName string, successRatio float64, interval time.Duration,
	serviceAdjuster *weight_cal.WeightAdjuster, zoneAdjuster *weight_cal.WeightAdjuster) {
	for i := 0; i < size; i++ {
		go func() {
			for range time.Tick(interval) {
				if rand.Float64() <= successRatio {
					serviceAdjuster.Notify(serviceName, balancer_common.Success)
					zoneAdjuster.Notify(zoneName, balancer_common.Success)
				} else {
					serviceAdjuster.Notify(serviceName, balancer_common.Failed)
					zoneAdjuster.Notify(zoneName, balancer_common.Failed)
				}
			}
		}()
	}
}

func NewServiceNode() []*api.ServiceEntry {
	service1 := &api.AgentService{
		Address: "192.168.1.1",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "local_zone", "__weight": "120"},
	}
	service2 := &api.AgentService{
		Address: "192.168.1.2",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "local_zone", "__weight": "100"},
	}
	service3 := &api.AgentService{
		Address: "192.168.1.3",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "local_zone", "__weight": "100"},
	}
	service4 := &api.AgentService{
		Address: "10.0.0.1",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "other_zone1", "__weight": "100"},
	}
	service5 := &api.AgentService{
		Address: "10.0.2.3",
		Port:    10000,
		Meta:    map[string]string{"__zone_id": "other_zone2", "__weight": "100"},
	}

	entrys := []*api.ServiceEntry{}
	entrys = append(entrys,
		&api.ServiceEntry{Service: service1},
		&api.ServiceEntry{Service: service2},
		&api.ServiceEntry{Service: service3},
		&api.ServiceEntry{Service: service4},
		&api.ServiceEntry{Service: service5})
	return entrys
}

type BalancerAdapter struct {
	balancer Balancer

	zoneAdjuster    *weight_cal.WeightAdjuster
	serviceAdjuster *weight_cal.WeightAdjuster
}

func (resolver *BalancerAdapter) UpdateServicesNotify(nodes []*balancer_common.ServiceNode) {
	//update weight_cal
	if resolver.balancer != nil {
		resolver.balancer.UpdateServices(nodes, resolver.zoneAdjuster, resolver.serviceAdjuster)
	}
}

func Test_WeightedRobinBalancer(t *testing.T) {
	Convey("Test_WeightedRobinBalancer", t, func() {
		//new Adjuster
		serviceAdjuster := weight_cal.NewWeightAdjuster()
		zoneAdjuster := weight_cal.NewWeightAdjuster()
		RandomNotify(50, "192.168.1.1:10000", "local_zone", 0.99, time.Duration(10)*time.Millisecond, serviceAdjuster, zoneAdjuster)
		RandomNotify(50, "192.168.1.2:10000", "local_zone", 0.25, time.Duration(10)*time.Millisecond, serviceAdjuster, zoneAdjuster)
		RandomNotify(50, "192.168.1.3:10000", "local_zone", 0.25, time.Duration(10)*time.Millisecond, serviceAdjuster, zoneAdjuster)
		RandomNotify(50, "10.0.0.1:10000", "other_zone1", 0.99, time.Duration(10)*time.Millisecond, serviceAdjuster, zoneAdjuster)
		RandomNotify(50, "10.0.2.3:10000", "other_zone2", 0.98, time.Duration(10)*time.Millisecond, serviceAdjuster, zoneAdjuster)
		//new balancer
		balancer := &WeightedRoundRobinBalancer{
			LocalZoneName: "local_zone",
		}
		//new adapter
		adapter := &BalancerAdapter{
			balancer:        balancer,
			zoneAdjuster:    zoneAdjuster,
			serviceAdjuster: serviceAdjuster,
		}
		//new discover
		discover := &discover.ConsulDiscover{}
		discover.UpdateNotify(adapter)

		entrys := NewServiceNode()
		discover.UpdateNodes(entrys)
		//fmt.Println("nodes:", balancer.Weights)
		//return
		go func() {
			for range time.Tick(time.Duration(5) * time.Second) {
				discover.UpdateNodes(entrys)
			}
		}()

		for j := 1; j <= 5; j++ {
			time.Sleep(time.Duration(1) * time.Second)
			countMap := make(map[string]int)
			for i := 0; i < 5000; i++ {
				time.Sleep(time.Duration(1) * time.Millisecond)
				node, err := balancer.DiscoverNode()
				if err == nil {
					countMap[node.Address] += 1
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
