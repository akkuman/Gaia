package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// ParseIPFormat parse ip block, for example: 192.168.0.1-192.168.5.8 192.168.4.5
func ParseIPFormat(ipStr string) (ips []string, err error) {
	if !strings.Contains(ipStr, "-") && !strings.Contains(ipStr, "/") {
		ips = append(ips, ipStr)
		return
	}
	if strings.Contains(ipStr, "-") && strings.Contains(ipStr, "/") {
		err = fmt.Errorf("%s format error", ipStr)
		return
	}
	// ip scope
	if strings.Contains(ipStr, "-") {
		ipStartEnd := strings.Split(ipStr, "-")
		ipStart := strings.TrimSpace(ipStartEnd[0])
		ipEnd := strings.TrimSpace(ipStartEnd[1])
		ips, err = generateIPListFromStart2End(ipStart, ipEnd)
	}
	// ip with mask
	if strings.Contains(ipStr, "/") {

	}
	return
}

func generateIPListFromStart2End(ipStart, ipEnd string) (ips []string, err error) {
	var startIPuint32 uint32
	var endIPuint32 uint32

	startIP := net.ParseIP(ipStart).To4()
	endIP := net.ParseIP(ipEnd).To4()

	startIPuint32 = binary.LittleEndian.Uint32(startIP)
	endIPuint32 = binary.LittleEndian.Uint32(endIP)

	for i := startIPuint32; i <= endIPuint32; i++ {
		var ipByte net.IP
		binary.LittleEndian.PutUint32(ipByte, i)
		ips = append(ips, ipByte.String())
	}
	return
}
