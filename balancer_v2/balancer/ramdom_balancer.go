package balancer

import (
	"errors"
	"math/rand"
	"sort"
	"sync"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

//struct RandomBalancer
type RandomBalancer struct {
	LocalZoneName string
	NodeName      string
	Weights       []*balancer_common.ServiceNode
	Factors       []int
	MaxFactors    int

	rwMutex sync.RWMutex
}

func NewRandomBalancer(localZoneName string, discoverNode string) Balancer {
	return &RandomBalancer{
		LocalZoneName: localZoneName,
		NodeName:      discoverNode,
	}
}

func (balancer *RandomBalancer) UpdateServices(nodes []*balancer_common.ServiceNode) error {
	factors := make([]int, 0, len(nodes))
	//cul Weight
	maxFactors := 0
	for _, node := range nodes {
		maxFactors += node.CurWeight
		//set nodes
		factors = append(factors, maxFactors)
	}
	balancer.rwMutex.Lock()
	defer balancer.rwMutex.Unlock()
	balancer.MaxFactors = maxFactors
	balancer.Factors = factors
	balancer.Weights = nodes
	return nil
}

func (balancer *RandomBalancer) DiscoverNode() (*balancer_common.ServiceNode, error) {
	size := int64(len(balancer.Factors))
	if size == 0 {
		return nil, errors.New("empty service nodes")
	}
	//lock
	balancer.rwMutex.RLock()
	maxFactors := balancer.MaxFactors
	factors := balancer.Factors
	nodes := balancer.Weights
	balancer.rwMutex.RUnlock()
	//get idx
	idx := sort.SearchInts(factors, rand.Intn(maxFactors))
	return nodes[idx], nil
}
