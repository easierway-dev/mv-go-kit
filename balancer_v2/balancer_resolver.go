package balancer_v2

import (
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/balancer"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/discover"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

type BalancerResolver struct {
	discover Discover
	balancer Balancer

	zoneAdjuster    *WeightAdjuster
	serviceAdjuster *WeightAdjuster

	nodes      []*ServiceNode //Service Node
	interval   time.Duration  //Update Time
	updateTime uint64

	mutex sync.Mutex
}

func NewBalancerResolver(balancerType, discoverType int, zoneName string,
	address string, discoverNode string, interval time.Duration) (*BalancerResolver, error) {
	//create resolver
	resolver = &BalancerResolver{interval: interval}
	//create balancer
	balancer, err := NewBalancer(balancerType, zoneName)
	if err != nil {
		return nil, err
	}
	resolver.balancer = balancer
	//create zone && service adjuster
	resolver.zoneAdjuster = weight_cal.NewWeightAdjuster()
	resolver.serviceAdjuster = weight_cal.NewWeightAdjuster()
	//create discover
	discover, err := NewDiscover(discoverType, address, interval, resolver)
	if err != nil {
		return nil, err
	}
	resolver.discover = discover
	//start UpdateBalancerByTimer
	resolver.UpdateBalancerByTimer(interval)
	return resolver
}

func (resolver *BalancerResolver) Notify(address string, zone string, event int) {
	if resolver.zoneAdjuster != nil {
		resolver.zoneAdjuster.Notify(zone, event)
	}
	if resolver.serviceAdjuster != nil {
		resolver.serviceAdjuster.Notify(address, event)
	}
}

func (resolver *BalancerResolver) UpdateBalancerByTimer(interval time.Duration) {
	go func() {
		for range time.Tick(interval) {
			resolver.UpdateServicesNotify(resolver.nodes)
		}
	}()
}

func (resolver *BalancerResolver) UpdateServicesNotify(nodes []*ServiceNode) {
	//lock
	resolver.mutex.Lock()
	defer resolver.mutex.Unlock()
	//updateTime
	now := time.Now().Unix()
	resolver.updateTime = time.Now().Unix()
	if now-resolver.updateTime < int64(resolver.interval) {
		return
	}
	//update weight_cal
	if resolver.balancer != nil {
		resolver.balancer.UpdateServices(nodes, resolver.zoneAdjuster, resolver.serviceAdjuster)
	}
}
