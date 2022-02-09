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

func (balancer *WeightedRoundRobinBalancer) UpdateServices(nodes []*balancer_common.ServiceNode,
	zoneAdjuster, serviceAdjuster *weight_cal.WeightAdjuster) error {
	weights := make([]*balancer_common.ServiceNode, 0, len(nodes))
	//open zone cul
	useZoneCul := false
	useZoneCulStr := "0"
	if len(balancer.LocalZoneName) != 0 {
		for _, node := range nodes {
			if balancer.LocalZoneName == node.Zone {
				useZoneCulStr = "1"
				useZoneCul = true
				break
			}
		}
	}
	//cul weight
	for _, node := range nodes {
		//cul zone Weight
		serviceWeight := weight_cal.GetServiceWeight(serviceAdjuster, node.Address)
		weight := float64(node.Weight) * serviceWeight
		zoneWeight := weight_cal.GetZoneWeight(zoneAdjuster, balancer.LocalZoneName, node.Zone)
		if useZoneCul {
			weight *= zoneWeight
		}
		culWeight := int(weight)
		//add metrics
		balancer_common.ZoneWeightHistogramVec.WithLabelValues(node.Zone, useZoneCulStr, balancer.NodeName).Observe(zoneWeight)
		balancer_common.IpWeightHistogramVec.WithLabelValues(node.Address, balancer.NodeName).Observe(serviceWeight)
		balancer_common.CulWeightHistogramVec.WithLabelValues(node.Address, node.Address, useZoneCulStr, balancer.NodeName).Observe(weight)
		//add node
		for i := 0; i < culWeight; i++ {
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
