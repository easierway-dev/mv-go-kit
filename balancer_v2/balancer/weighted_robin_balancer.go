package balancer

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

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
	zoneAdjuster, serviceAdjuster *weight_cal.WeightAdjuster) error {
	balancer.Weights = make([]*balancer_common.ServiceNode, 0, len(nodes))
	for _, node := range nodes {
		//cul zone Weight
		weight := float64(node.Weight) *
			weight_cal.GetZoneWeight(zoneAdjuster, balancer.LocalZoneName, node.Zone) *
			weight_cal.GetServiceWeight(serviceAdjuster, node.Zone)
		culWeight := int(weight)
		if culWeight <= 0 {
			culWeight = 1
		}
		//add node
		for i := 0; i < culWeight; i++ {
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
	size := int64(len(balancer.Weights))
	if size == 0 {
		return nil, errors.New("empty service nodes")
	}
	idx := int(atomic.AddInt64(&balancer.Count, 1) % size)
	return balancer.Weights[idx], nil
}
