package discover

import (
	"testing"

	"github.com/hashicorp/consul/api"

	. "github.com/smartystreets/goconvey/convey"
)

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

func Test_UpdateNodes(t *testing.T) {
	Convey("Test_UpdateNodes", t, func() {
		discover := &ConsulDiscover{}
		entrys := NewServiceNode()
		discover.UpdateNodes(entrys)
		So(len(discover.nodes), ShouldEqual, 5)
	})
}
