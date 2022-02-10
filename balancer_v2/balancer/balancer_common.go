package balancer

import (
	"gitlab.mobvista.com/voyager/mv-go-kit/balancer_v2/common"
)

func CheckOpenZoneWeight(nodes []*balancer_common.ServiceNode, localZoneName string) bool {
	localZoneNum := 0
	otherZoneNum := 0
	if len(localZoneName) != 0 {
		for _, node := range nodes {
			if localZoneName == node.Zone {
				localZoneNum += 1
			} else {
				otherZoneNum += 1
			}
		}
	}
	if localZoneNum > 0 && otherZoneNum > 0 {
		return true
	} else {
		return false
	}
}
