package plugin

import (
	"fmt"
	"sync"
)

// ThreadFunc thread execution cell
type ThreadFunc func(info [3]string) (success bool, err error)

// SuccessMap success map
var SuccessMap = make(map[string][3]string)

// SwitchIsOn Whether the plugin is enabled
func SwitchIsOn(pluginName string) bool {
	if _, ok := OnOptions[pluginName]; ok {
		return true
	}
	return false
}

// TaskCellStart start a task cell
func TaskCellStart(pluginName string, wgMain *sync.WaitGroup, threadFunc ThreadFunc) {
	wg := new(sync.WaitGroup)
	tasks := make(chan [3]string, 100)
	wg.Add(ThreadNum)

	for i := 0; i < ThreadNum; i++ {
		go thread(pluginName, tasks, wg, i, threadFunc)
	}

	var burstCell = NewBurstCell(pluginName)

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

func thread(pluginName string, tasks chan [3]string, wg *sync.WaitGroup, threadNo int, threadfunc ThreadFunc) {
	for {
		cell, ok := <-tasks
		if !ok {
			break
		}
		if _, ok = SuccessMap[cell[0]]; ok {
			continue
		}
		CliOutput.LastPath(fmt.Sprintf("[!]Start: %s %s %s %s", pluginName, cell[0], cell[1], cell[2]), threadNo)
		success, err := threadfunc(cell)
		if err != nil {
			CliOutput.LastPath(fmt.Sprintf("[x]Fail: %s %s", pluginName, err.Error()), threadNo)
			continue
		}
		if success {
			CliOutput.StatusReport(fmt.Sprintf("[*]Find: %6s %21s %16s %16s", pluginName, cell[0], cell[1], cell[2]))
			SuccessMap[cell[0]] = [3]string{pluginName, cell[1], cell[2]}
		}
	}
	wg.Done()
}
