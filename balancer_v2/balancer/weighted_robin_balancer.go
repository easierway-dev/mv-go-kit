package balancer

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

type WeightedRoundRobinBalancer struct {
	LocalZoneName string
	NodeName      string
	Count         int64
	Weights       []*balancer_common.ServiceNode
}

func NewWeightedRoundRobin(localZoneName string, discoverNode string) Balancer {
	return &WeightedRoundRobinBalancer{
		LocalZoneName: localZoneName,
		NodeName:      discoverNode,
	}
}

func (balancer *WeightedRoundRobinBalancer) UpdateServices(nodes []*balancer_common.ServiceNode) error {
	weights := make([]*balancer_common.ServiceNode, 0, len(nodes)*50)
	//cul weight
	for _, node := range nodes {
		//cul zone Weight
		curWeight := node.CurWeight
		//add node
		for i := 0; i < curWeight; i++ {
			weights = append(weights, node)
		}
	}
	//shuffle node
	balancer.Count = 0
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(weights), func(i, j int) {
		weights[i], weights[j] = weights[j], weights[i]
	})
	balancer.Weights = weights
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
