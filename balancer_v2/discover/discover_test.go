package discover

import (
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/hashicorp/consul/api"
	. "github.com/smartystreets/goconvey/convey"

	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

func Test_ConsulDiscover(t *testing.T) {
	Convey("Test_ConsulDiscover", t, func() {
		patches := ApplyFunc(discover.client.Health().Service, func(string, string, bool, *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {

		})
		patches.Reset()
	})
}
