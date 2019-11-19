package ssh

import (
	"Gaia/plugin"
	"Gaia/util"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

var pluginName = "ssh"

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
	config := &ssh.ClientConfig{
		User: info[1],
		Auth: []ssh.AuthMethod{
			ssh.Password(info[2]),
		},
		Timeout: util.TimeOut,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", info[0], config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo xsec")
		if err == nil && errRet == nil {
			defer session.Close()
			success = true
		}
	}
	return
}

func init() {
	ftpPlugin := BurstPlugin{}
	plugin.Register(pluginName, ftpPlugin)
}
