package discover

import (
	"errors"
	"time"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

type DiscoverNotify interface {
	UpdateServicesNotify(nodes []*balancer_common.ServiceNode)
}

type Discover interface {
	Start() error
}

func NewDiscover(discoverType int, address string, discoverNode string,
	interval time.Duration, notify DiscoverNotify, logger balancer_common.Logger) (Discover, error) {
	switch discoverType {
	case balancer_common.ConsulDiscover:
		discover, err := NewConsulDiscover(address, discoverNode, interval, notify, logger)
		if err != nil {
			return nil, err
		}
		return discover, nil
	default:
		return nil, errors.New("undefine discover type")
	}
}
