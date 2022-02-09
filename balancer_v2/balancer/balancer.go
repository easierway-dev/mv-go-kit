package balancer

import (
	"errors"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

// Balancer interface
type Balancer interface {
	DiscoverNode() (*balancer_common.ServiceNode, error)
	UpdateServices(nodes []*balancer_common.ServiceNode, zoneAdjuster, serviceAdjuster *weight_cal.WeightAdjuster) error
}

// BalancerSelector
type BalancerSelector struct {
	balancerType int
	balancer     Balancer
}

func NewBalancer(balancerType int, zoneName string, discoverNode string) (Balancer, error) {
	switch balancerType {
	case balancer_common.RoundRobin, balancer_common.ConsistentHash:
		return nil, errors.New("no support")
	case balancer_common.RandomSelect:
		return NewRandomBalancer(zoneName, discoverNode), nil
	case balancer_common.WeightedRoundRobin:
		return NewWeightedRoundRobin(zoneName, discoverNode), nil
	default:
		return nil, errors.New("undefine balance type")
	}
}
