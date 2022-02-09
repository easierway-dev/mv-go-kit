package discover

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

// Test Discover Notify
type MyDiscoverNotify struct {
}

func (notify *MyDiscoverNotify) UpdateServicesNotify(nodes []*balancer_common.ServiceNode) {
	for _, node := range nodes {
		fmt.Println("notify nodes:", node)
	}
	fmt.Println("")
}

//test logger
type MyLogger struct {
}

func (logger *MyLogger) Infof(format string, v ...interface{}) {}
func (logger *MyLogger) Warnf(format string, v ...interface{}) {}

func NewServiceNode() []*api.ServiceEntry {
	service1 := &api.AgentService{
		Address: "192.168.1.1",
		Port:    9888,
		Meta:    map[string]string{"__zone_id": "A", "__weight": "120"},
	}
	service2 := &api.AgentService{
		Address: "192.168.1.2",
		Port:    9888,
		Meta:    map[string]string{"__zone_id": "A", "__weight": "100"},
	}
	service3 := &api.AgentService{
		Address: "192.168.1.3",
		Port:    9888,
		Meta:    map[string]string{"__zone_id": "A", "__weight": "100"},
	}
	service4 := &api.AgentService{
		Address: "10.10.1.2",
		Port:    9888,
		Meta:    map[string]string{"__zone_id": "B", "__weight": "100"},
	}
	service5 := &api.AgentService{
		Address: "10.10.2.2",
		Port:    9888,
		Meta:    map[string]string{"__zone_id": "C", "__weight": "100"},
	}

	entrys := []*api.ServiceEntry{}
	entrys = append(entrys,
		&api.ServiceEntry{Service: service1},
		&api.ServiceEntry{Service: service2},
		&api.ServiceEntry{Service: service3},
		&api.ServiceEntry{Service: service4},
		&api.ServiceEntry{Service: service5})
	return entrys
}

func Test_ConsulDiscover(t *testing.T) {
	Convey("Test_ConsulDiscover", t, func() {
		discover := &ConsulDiscover{
			notify: &MyDiscoverNotify{},
		}
		entrys := NewServiceNode()
		go func() {
			for range time.Tick(time.Duration(2) * time.Second) {
				discover.UpdateNodes(entrys)
			}
		}()
		time.Sleep(time.Duration(10) * time.Second)
	})
}
