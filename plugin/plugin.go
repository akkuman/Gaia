package plugin

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"Gaia/util"

	"gopkg.in/ini.v1"
)

// IPs a space with ip list
var IPs []string

// OnOptions which service will be burst
var OnOptions map[string]bool

// Config plugin config
type Config struct {
	IPStr   string
	options []string // which service will be burst
}

// Plugins a container to store plugin
var Plugins map[string]Plugin

// Plugin plugin interface
type Plugin interface {
	Flag() bool
	Start()
}

// BurstCell a burst task cell
type BurstCell struct {
	Usernames []string
	Passwords []string
	Ports     []int
	IPs       []string
}

// InitConfig initialize config
func initConfig(config Config) {
	var err error
	IPs, err = util.ParseIPFormat(config.IPStr)
	for _, v := range config.options {
		OnOptions[v] = true
	}
	if err != nil {
		panic(err)
	}
}

// Start start all plugin from plugin container
func (config Config) Start() {
	initConfig(config)
	for name, plugin := range Plugins {
		if plugin.Flag() {
			go plugin.Start()
			fmt.Printf("[*]Start plugin %s\n", name)
		} else {
			fmt.Printf("[-]Skip plugin %s", name)
		}
	}
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
	if len(IPs) == 0 {
		panic("ip is uninitialized")
	}
	burstCell.IPs = IPs
	return
}

// Register register a plugin to plugin container
func Register(name string, plugin Plugin) {
	Plugins[name] = plugin
}
