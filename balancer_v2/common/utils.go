package balancer_common

import (
	"errors"
	"net"
	"regexp"
)

var ip4Reg = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

func GetLocalIp() (string, error) {
	addr, err := localIPv4s()
	if err != nil {
		return "", err
	}

	if len(addr) == 0 {
		return "", errors.New("get local ip error")
	}
	return addr[0], nil
}

func localIPv4s() ([]string, error) {
	var ips []string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, a := range addr {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			if ip4Reg.MatchString(ipnet.IP.String()) {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

func CheckOpenZoneWeight(nodes []*ServiceNode, localZoneName string) bool {
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
