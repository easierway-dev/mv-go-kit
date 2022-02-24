package resource

import (
	"log"
	"net"
	//   "strings"
)

var LocalIP string

func init() {
	LocalIP = GetOutboundIP()
}
func IP() string {
	return LocalIP
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
