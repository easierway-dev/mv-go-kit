package balancer

import (
	"sort"
	"sync/atomic"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

type WeightedRoundRobinBalancer struct {
	LocalZoneName string
	Count         int64
	Weights       []*balancer_common.ServiceNode
}

func NewWeightedRoundRobin(localZoneName string) Balancer {
	return &WeightedRoundRobinBalancer{
		LocalZoneName: localZoneName,
	}
}

func (balancer *WeightedRoundRobinBalancer) UpdateServices(nodes []*balancer_common.ServiceNode,
	zoneAdjuster, serviceAdjuster *weight_cal.WeightAdjuster) {
	balancer.Weights = make([]*balancer_common.ServiceNode, 0, len(nodes))
	for _, node := range nodes {
		//cul zone Weight
		weight := node.Weight *
			weight_cal.GetZoneWeight(zoneAdjuster, balancer.LocalZoneName, node.Zone) *
			weight_cal.GetServiceWeight(serviceAdjuster, node.Zone)
		//add node
		for i := 0; i < weight; i++ {
			balancer.Weights = append(balancer.Weights, node)
		}
	}
	//shuffle node
	balancer.Count = 0
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(balancer.Weights), func(i, j int) {
		balancer.Weights[i], balancer.Weights[j] = balancer.Weights[j], balancer.Weights[i]
	})
	return nil
}

func (balancer *WeightedRoundRobinBalancer) DiscoverNode() (*balancer_common.ServiceNode, error) {
	if len(balancer.Weights) == 0 {
		return nil, errors.New("empty service nodes")
	}
	idx := atomic.AddInt64(balancer.Count, 1) % len(balancer.Weights)
	return balancer.Weights[idx], nil
}
