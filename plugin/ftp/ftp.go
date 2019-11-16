package ftp

import (
	"Gaia/plugin"
)

var pluginName = "ftp"

var burstCell = plugin.NewBurstCell(pluginName)

// BurstPlugin ftp burst plugin
type BurstPlugin struct {
}

// Flag turn on?
func (BurstPlugin) Flag() bool {
	if _, ok := plugin.OnOptions[pluginName]; ok {
		return true
	}
	return false
}

// Start start a burst
func (BurstPlugin) Start() {

}

func init() {
	ftpPlugin := BurstPlugin{}
	plugin.Register(pluginName, ftpPlugin)
}
