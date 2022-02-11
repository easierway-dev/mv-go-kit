package balancer_v2

import (
	"time"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/balancer"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/discover"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

type BalancerResolver struct {
	discover discover.Discover //service discover
	balancer balancer.Balancer //balancer

	discoverNode string
	localZone    string

	zoneAdjuster    *weight_cal.WeightAdjuster
	serviceAdjuster *weight_cal.WeightAdjuster
}

func NewBalancerResolver(balancerType, discoverType int, zoneName string,
	address string, discoverNode string, interval time.Duration, logger balancer_common.Logger) (*BalancerResolver, error) {
	//create resolver
	resolver := &BalancerResolver{}
	//create balancer
	balancer, err := balancer.NewBalancer(balancerType, zoneName, discoverNode)
	if err != nil {
		return nil, err
	}
	resolver.discoverNode = discoverNode
	resolver.localZone = zoneName
	resolver.balancer = balancer
	//create zone && service adjuster
	resolver.zoneAdjuster = weight_cal.NewWeightAdjuster()
	resolver.serviceAdjuster = weight_cal.NewWeightAdjuster()
	//create discover
	discover, err := discover.NewDiscover(discoverType, address, discoverNode, interval, resolver, logger)
	if err != nil {
		return nil, err
	}
	resolver.discover = discover
	return resolver, nil
}

func (resolver *BalancerResolver) Notify(address string, zone string, event int) {
	if resolver.zoneAdjuster != nil {
		resolver.zoneAdjuster.Notify(zone, event)
	}
	if resolver.serviceAdjuster != nil {
		resolver.serviceAdjuster.Notify(address, event)
	}
}

func (resolver *BalancerResolver) UpdateServicesNotify(nodes []*balancer_common.ServiceNode) {
	//open zone cul
	useZoneCul := CheckOpenZoneWeight(nodes, resolver.localZone)
	useZoneCulStr := "0"
	if useZoneCul {
		useZoneCulStr = "1"
	}
	//cal CurWeight
	for _, node := range nodes {
		//cul zone Weight
		serviceWeight := weight_cal.GetServiceWeight(resolver.serviceAdjuster, node.Address)
		weight := float64(node.Weight) * serviceWeight
		zoneWeight := weight_cal.GetZoneWeight(resolver.zoneAdjuster, resolver.localZone, node.Zone)
		if useZoneCul {
			weight *= zoneWeight
		}
		node.CurWeight = int(weight)
		//local zone min weight = 1
		if node.CurWeight == 0 && node.Zone == resolver.localZone {
			node.CurWeight = 1
		}
		//add metrics
		balancer_common.ZoneWeightHistogramVec.WithLabelValues(node.Zone, useZoneCulStr, resolver.discoverNode).Observe(zoneWeight)
		balancer_common.IpWeightHistogramVec.WithLabelValues(node.Address, resolver.discoverNode).Observe(serviceWeight)
		balancer_common.CulWeightHistogramVec.WithLabelValues(node.Address, node.Address, useZoneCulStr, resolver.discoverNode).Observe(weight)
	}
	//update weight_cal
	if resolver.balancer != nil {
		resolver.balancer.UpdateServices(nodes)
	}
}

func (resolver *BalancerResolver) DiscoverNode() (*balancer_common.ServiceNode, error) {
	return resolver.balancer.DiscoverNode()
}

func (resolver *BalancerResolver) GetNode() (string, error) {
	node, err := resolver.DiscoverNode()
	if err != nil {
		return "", err
	}
	//add metrics
	balancer_common.ZoneIpCallCounter.WithLabelValues(node.Zone, node.Address, resolver.discoverNode).Inc()
	return node.Address, nil
}

func CheckOpenZoneWeight(nodes []*balancer_common.ServiceNode, localZoneName string) bool {
	localZoneNum := 0
	otherZoneNum := 0
	if len(localZoneName) != 0 {
		for _, node := range nodes {
			if localZoneName == node.Zone {
				localZoneNum += 1
			} else {
				otherZoneNum += 1
			}
		}
	}
	if localZoneNum > 0 && otherZoneNum > 0 {
		return true
	} else {
		return false
	}
}
