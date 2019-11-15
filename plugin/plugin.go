package plugin

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"Gaia/util"

	"gopkg.in/ini.v1"
)

type PluginInfo struct {
	IPStr string
}

var IPs []string

// BurstFunc function prototype for burst
type BurstFunc func([]string, []string)

// BurstCell a burst task cell
type BurstCell struct {
	Usernames []string
	Passwords []string
	Ports     []int
	IPs       []string
}

// NewBurstCell generate a new burst task cell
func NewBurstCell(burstType string) (burstCell BurstCell) {
	// 读取配置
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	usernameStr := cfg.Section(burstType).Key("username").String()
	passwordStr := cfg.Section(burstType).Key("password").String()
	PortStr := cfg.Section(burstType).Key("ports").String()

	burstCell.Usernames, err = util.GetConfigStrList(usernameStr)
	if err != nil {
		panic(err)
	}
	burstCell.Passwords, err = util.GetConfigStrList(passwordStr)
	if err != nil {
		panic(err)
	}
	for _, v := range strings.Split(PortStr, ",") {
		item := strings.TrimSpace(v)
		port, err := strconv.Atoi(item)
		if err != nil {
			panic(err)
		}
		burstCell.Ports = append(burstCell.Ports, port)
	}
	return
}

// InitPluginInfo init plugin info
func InitPluginInfo(pluginInfo PluginInfo) (err error) {
	for _, ipBlock := range strings.Split(pluginInfo.IPStr, ",") {
		// ...
	}
	return
}

// Run start a burst
func (BurstCell) Run()
