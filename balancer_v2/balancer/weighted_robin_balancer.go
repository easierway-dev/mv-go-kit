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

//辗转相除法求最大公约数
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func GetGcd(nodes []*balancer_common.ServiceNode) int {
	if len(nodes) == 0 {
		return 0
	}
	g := nodes[0].CurWeight
	for _, node := range nodes {
		curWeight := node.CurWeight
		if curWeight == 0 {
			continue
		}
		//oldGcd := g
		g = gcd(g, curWeight)
	}
	return g
}

func GetWeightCount(nodes []*balancer_common.ServiceNode) int {
	count := 0
	for _, node := range nodes {
		count += node.CurWeight
	}
	return count
}

func (balancer *WeightedRoundRobinBalancer) UpdateServices(nodes []*balancer_common.ServiceNode) error {
	gcd := GetGcd(nodes)
	weights := make([]*balancer_common.ServiceNode, 0, GetWeightCount(nodes))
	//cul weight
	for _, node := range nodes {
		//cul zone Weight
		curWeight := node.CurWeight
		if gcd != 0 {
			curWeight /= gcd
		}
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
