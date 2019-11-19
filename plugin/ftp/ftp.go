package ftp

import (
	"Gaia/plugin"
	"Gaia/util"
	"sync"

	"github.com/jlaffaye/ftp"
)

var pluginName = "ftp"

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
	c, err := ftp.Dial(info[0], ftp.DialWithTimeout(util.TimeOut))
	if err != nil {
		return
	}

	err = c.Login(info[1], info[2])
	if err != nil {
		return
	}
	success = true
	_ = c.Quit()
	return
}

func init() {
	ftpPlugin := BurstPlugin{}
	plugin.Register(pluginName, ftpPlugin)
}
