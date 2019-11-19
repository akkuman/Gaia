package plugin

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"Gaia/clioutput"
	"Gaia/util"

	"gopkg.in/ini.v1"
)

// IPs a space with ip list
var IPs []string

// OnOptions which service will be burst
var OnOptions = make(map[string]bool)

// ThreadNum thread num
var ThreadNum int

// CliOutput cli output
var CliOutput = clioutput.NewCLIOutput()

// Config plugin config
type Config struct {
	IPStr     string
	Options   []string // which service will be burst
	ThreadNum int
}

// Plugins a container to store plugin
var Plugins = make(map[string]Plugin)

// Plugin plugin interface
type Plugin interface {
	Flag() bool
	Start(*sync.WaitGroup)
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
	for _, v := range config.Options {
		OnOptions[v] = true
	}
	if err != nil {
		panic(err)
	}
	ThreadNum = config.ThreadNum
}

// Start start all plugin from plugin container
func (config Config) Start() {
	initConfig(config)
	wg := new(sync.WaitGroup)
	for name, plugin := range Plugins {
		if plugin.Flag() {
			wg.Add(1)
			CliOutput.StatusReport(fmt.Sprintf("[*]Start plugin %s", name))
			go plugin.Start(wg)
		} else {
			CliOutput.StatusReport(fmt.Sprintf("[-]Skip plugin %s", name))
		}
	}
	wg.Wait()
}

// NewBurstCell generate a new burst task cell
func NewBurstCell(burstType string) (burstCell BurstCell) {
	// 读取配置
	cfg, err := ini.Load("config.ini")
	if err != nil {
		CliOutput.StatusReport(fmt.Sprintf("Fail to read file: %v", err))
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
