package ftp

import (
	"Gaia/plugin"
	"fmt"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

var pluginName = "ftp"

var successMap = make(map[string]bool)

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
func (BurstPlugin) Start(wgMain *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	tasks := make(chan [3]string, 100)
	wg.Add(plugin.ThreadNum)

	for i := 0; i < plugin.ThreadNum; i++ {
		go thread(tasks, wg, i)
	}

	var burstCell = plugin.NewBurstCell(pluginName)

	for _, ip := range burstCell.IPs {
		for _, port := range burstCell.Ports {
			for _, username := range burstCell.Usernames {
				for _, password := range burstCell.Passwords {
					tasks <- [3]string{fmt.Sprintf("%s:%d", ip, port), username, password}
				}
			}
		}
	}
	close(tasks)
	wg.Wait()
	wgMain.Done()
}

func thread(tasks chan [3]string, wg *sync.WaitGroup, threadNo int) {
	for {
		cell, ok := <-tasks
		if !ok {
			break
		}
		if _, ok = successMap[cell[0]]; ok {
			continue
		}
		plugin.CliOutput.LastPath(fmt.Sprintf("[!]Start: ftp %s %s %s", cell[0], cell[1], cell[2]), threadNo)
		success, err := threadFunc(cell)
		if err != nil {
			plugin.CliOutput.LastPath(fmt.Sprintf("[x]Fail: ftp %s", err.Error()), threadNo)
			continue
		}
		if success {
			plugin.CliOutput.StatusReport(fmt.Sprintf("[*]Find: ftp %s %s %s", cell[0], cell[1], cell[2]))
			successMap[cell[0]] = true
		}
	}
	wg.Done()
}

func threadFunc(info [3]string) (success bool, err error) {
	c, err := ftp.Dial(info[0], ftp.DialWithTimeout(5*time.Second))
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
