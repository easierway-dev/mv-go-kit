package balancer

import (
	"errors"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/weight_cal"
)

//struct RandomBalancer
type RandomBalancer struct {
	LocalZoneName string
	NodeName      string
}

func NewRandomBalancer(localZoneName string, discoverNode string) Balancer {
	return &RandomBalancer{
		LocalZoneName: localZoneName,
		NodeName:      discoverNode,
	}
}

func (balancer *RandomBalancer) UpdateServices(nodes []*balancer_common.ServiceNode, zoneAdjuster, serviceAdjuster *weight_cal.WeightAdjuster) error {
	return nil
}

func (balancer *RandomBalancer) DiscoverNode() (*balancer_common.ServiceNode, error) {
	return nil, errors.New("not support")
}
