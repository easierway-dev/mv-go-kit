package discover

import (
	"errors"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

type DiscoverNotify interface {
	UpdateServicesNotify(nodes []*balancer_common.ServiceNode)
}

type Discover interface {
	Start() error
}

func NewDiscover(discoverType int, address string, discoverNode string,
	interval time.Duration, notify DiscoverNotify) (Discover, error) {
	switch discoverType {
	case ConsulDiscover:
		discover, err := NewConsulDiscover(address, discoverNode, interval, notify)
		if err != nil {
			return nil, err
		}
		return discover, nil
	default:
		return nil, errors.New("undefine discover type")
	}
}
