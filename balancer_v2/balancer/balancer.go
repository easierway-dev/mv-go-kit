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

func NewBalancer(balancerType int, zoneName string) (Balancer, error) {
	switch balancerType {
	case RoundRobin, ConsistentHash:
		return nil, errors.New("no support")
	case RandomSelect:
		return NewRandomBalancer(zoneName)
	case WeightedRoundRobin:
		return NewWeightedRoundRobin(zoneName)
	default:
		return nil, errors.New("undefine balance type")
	}
}
