package ssh

import (
	"Gaia/plugin"
	"Gaia/util"
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
		Timeout:         util.TimeOut,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", info[0], config)
	if err != nil {
		return
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	success = true
	return
}

func init() {
	sshPlugin := BurstPlugin{}
	plugin.Register(pluginName, sshPlugin)
}
