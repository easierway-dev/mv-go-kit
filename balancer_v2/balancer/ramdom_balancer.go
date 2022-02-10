package balancer

import (
	"errors"
	"math/rand"
	"sort"
	"sync"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
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

func (balancer *RandomBalancer) UpdateServices(nodes []*balancer_common.ServiceNode, zoneAdjuster, serviceAdjuster *weight_cal.WeightAdjuster) error {
	factors := make([]int, 0, len(nodes))
	//open zone cul
	useZoneCul := CheckOpenZoneWeight(nodes, balancer.LocalZoneName)
	useZoneCulStr := "0"
	if useZoneCul {
		useZoneCulStr = "1"
	}
	//cul Weight
	maxFactors := 0
	for _, node := range nodes {
		//cul zone Weight
		serviceWeight := weight_cal.GetServiceWeight(serviceAdjuster, node.Address)
		weight := float64(node.Weight) * serviceWeight
		zoneWeight := weight_cal.GetZoneWeight(zoneAdjuster, balancer.LocalZoneName, node.Zone)
		if useZoneCul {
			weight *= zoneWeight
		}
		culWeight := int(weight)
		maxFactors += culWeight
		//add metrics
		balancer_common.ZoneWeightHistogramVec.WithLabelValues(node.Zone, useZoneCulStr, balancer.NodeName).Observe(zoneWeight)
		balancer_common.IpWeightHistogramVec.WithLabelValues(node.Address, balancer.NodeName).Observe(serviceWeight)
		balancer_common.CulWeightHistogramVec.WithLabelValues(node.Address, node.Address, useZoneCulStr, balancer.NodeName).Observe(weight)
		//set nodes
		factors = append(factors, maxFactors)
	}
	balancer.rwMutex.Lock()
	balancer.MaxFactors = maxFactors
	balancer.Factors = factors
	balancer.Weights = nodes
	balancer.rwMutex.Unlock()
	return nil
}

func (balancer *RandomBalancer) DiscoverNode() (*balancer_common.ServiceNode, error) {
	size := int64(len(balancer.Weights))
	if size == 0 {
		return nil, errors.New("empty service nodes")
	}
	//lock
	balancer.rwMutex.RLock()
	maxFactors := balancer.MaxFactors
	factors := balancer.Factors
	nodes := balancer.Weights
	balancer.rwMutex.RUnlock()
	idx := sort.SearchInts(factors, rand.Intn(maxFactors))
	return nodes[idx], nil
}
