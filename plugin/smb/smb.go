package smb

import (
	"Gaia/plugin"
	"strconv"
	"strings"
	"sync"

	"github.com/stacktitan/smb/smb"
)

var pluginName = "smb"

// BurstPlugin ftp burst plugin
type BurstPlugin struct {
}

// Flag turn on?
func (BurstPlugin) Flag() bool {
	return plugin.SwitchIsOn(pluginName)
}

// Start start a burst
func (BurstPlugin) Start(wgMain *sync.WaitGroup) {
	plugin.TaskCellStart(pluginName, wgMain, threadFunc)
}

func threadFunc(info [3]string) (success bool, err error) {
	ipport := strings.Split(info[0], ":")
	port, err := strconv.Atoi(ipport[1])
	if err != nil {
		return
	}
	options := smb.Options{
		Host:        ipport[0],
		Port:        port,
		User:        info[1],
		Password:    info[2],
		Domain:      "",
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err != nil {
		return
	}
	defer session.Close()
	if session.IsAuthenticated {
		success = true
	}
	return
}

func init() {
	smbPlugin := BurstPlugin{}
	plugin.Register(pluginName, smbPlugin)
}
