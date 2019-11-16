package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
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
		if len(ipStartEnd) != 2 {
			err = fmt.Errorf("%s format error", ipStr)
			return
		}
		ipStart := strings.TrimSpace(ipStartEnd[0])
		ipEnd := strings.TrimSpace(ipStartEnd[1])
		ips, err = generateIPListFromStart2End(ipStart, ipEnd)
	}
	// ip with mask
	if strings.Contains(ipStr, "/") {
		ipMask := strings.Split(ipStr, "/")
		if len(ipMask) != 2 {
			err = fmt.Errorf("%s format error", ipStr)
			return
		}
		ip := strings.TrimSpace(ipMask[0])
		mask := strings.TrimSpace(ipMask[1])
		ips, err = generateIPListViaIPMask(ip, mask)
	}
	return
}

// generateIPListFromStart2End generate ip list from ip scope
func generateIPListFromStart2End(ipStart, ipEnd string) (ips []string, err error) {
	var startIPuint32 uint32
	var endIPuint32 uint32

	startIP := net.ParseIP(ipStart).To4()
	endIP := net.ParseIP(ipEnd).To4()

	startIPuint32 = binary.BigEndian.Uint32(startIP)
	endIPuint32 = binary.BigEndian.Uint32(endIP)

	for i := startIPuint32; i <= endIPuint32; i++ {
		ipByte := make([]byte, 4)
		binary.BigEndian.PutUint32(ipByte, i)
		ips = append(ips, net.IP(ipByte).String())
	}
	return
}

// generateIPListViaIPMask generate ip list form ip with mask
func generateIPListViaIPMask(ip, mask string) (ips []string, err error) {
	var temp0 uint32
	var temp1 uint32
	baseIP := net.ParseIP(ip).To4()
	ipMask, err := strconv.Atoi(mask)
	if err != nil {
		return
	}
	if ipMask > 32 {
		err = fmt.Errorf("ip mask is out of range, It's 32 at most")
		return
	}
	temp0 = 0
	temp1 = 0xffffffff
	baseIPuint32 := binary.BigEndian.Uint32(baseIP)
	startIPuint32 := uint32(uint32(baseIPuint32>>(32-ipMask))<<(32-ipMask)) | uint32(temp0)       // 先右移去除右边的，再左移
	endIPuint32 := uint32(uint32(baseIPuint32>>(32-ipMask))<<(32-ipMask)) | uint32(temp1>>ipMask) // 先右移去除右边的，再左移
	for i := startIPuint32; i <= endIPuint32; i++ {
		ipByte := make([]byte, 4)
		binary.BigEndian.PutUint32(ipByte, i)
		ips = append(ips, net.IP(ipByte).String())
	}
	return
}
