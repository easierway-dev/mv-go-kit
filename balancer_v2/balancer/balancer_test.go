package balancer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

func NewDiscoverNode() []*balancer_common.ServiceNode {
	node1 := &balancer_common.ServiceNode{
		Address:   "192.168.1.1:10000",
		Host:      "192.168.1.1",
		Port:      10000,
		Zone:      "local_zone",
		Weight:    120,
		CurWeight: 120,
	}
	node2 := &balancer_common.ServiceNode{
		Address:   "192.168.1.2:10000",
		Host:      "192.168.1.2",
		Port:      10000,
		Zone:      "local_zone",
		Weight:    100,
		CurWeight: 100,
	}
	node3 := &balancer_common.ServiceNode{
		Address:   "192.168.1.3:10000",
		Host:      "192.168.1.3",
		Port:      10000,
		Zone:      "local_zone",
		Weight:    100,
		CurWeight: 100,
	}
	node4 := &balancer_common.ServiceNode{
		Address:   "10.0.0.1:10000",
		Host:      "10.0.0.1",
		Port:      10000,
		Zone:      "other_zone1",
		Weight:    100,
		CurWeight: 100,
	}
	node5 := &balancer_common.ServiceNode{
		Address:   "10.0.2.3:10000",
		Host:      "10.0.2.3",
		Port:      10000,
		Zone:      "other_zone2",
		Weight:    100,
		CurWeight: 100,
	}
	//add nodes
	nodes := make([]*balancer_common.ServiceNode, 0, 5)
	nodes = append(nodes, node1, node2, node3, node4, node5)
	return nodes
}

func Test_RamdomBalancer(t *testing.T) {
	Convey("Test_RamdomBalancer", t, func() {
		//init RandomBalancer
		balancer := &RandomBalancer{
			LocalZoneName: "local_zone",
			NodeName:      "test_node",
		}
		//create nodes
		nodes := NewDiscoverNode()
		//update nodes
		balancer.UpdateServices(nodes)
		//check factors
		So(len(balancer.Factors), ShouldEqual, 5)
		//discover nodes
		_, err := balancer.DiscoverNode()
		//check err
		So(err, ShouldEqual, nil)
	})
}

func Test_WeightedRoundRobinBalancer(t *testing.T) {
	Convey("Test_WeightedRoundRobinBalancer", t, func() {
		//init RandomBalancer
		balancer := &WeightedRoundRobinBalancer{
			LocalZoneName: "local_zone",
			NodeName:      "test_node",
		}
		//create nodes
		nodes := NewDiscoverNode()
		//update nodes
		balancer.UpdateServices(nodes)
		//check weight
		So(len(balancer.Weights), ShouldEqual, 26)
		//discover nodes
		_, err := balancer.DiscoverNode()
		//check err
		So(err, ShouldEqual, nil)
	})
}

func Test_GCD(t *testing.T) {
	Convey("Test_WeightedRoundRobinBalancer", t, func() {
		// create nodes
		nodes := NewDiscoverNode()
		//get gcd
		g := GetGcd(nodes)
		//check
		So(g, ShouldEqual, 20)
	})
}
