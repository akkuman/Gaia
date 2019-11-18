package clioutput

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

// CLIOuotput 控制台的一些信息
type CLIOuotput struct {
	lastLength int
	lastOutput string
	lastInLine bool
	mutex      *sync.Mutex
}

// NewCLIOutput new一个CLIOutput结构体
func NewCLIOutput() CLIOuotput {
	clioutpt := initCLIOuotput()
	return clioutpt
}

// LastPath 底行刷新
// func (cli *CLIOuotput) LastPath(info string, index, length int) {
// 	cli.mutex.Lock()
// 	defer cli.mutex.Unlock()

// 	x, _ := getTerminalSize()

// 	message := fmt.Sprintf("%04.2f%% - ", percentage(index, length))
// 	message = fmt.Sprintf("%s%s", message, info)

// 	if len(message) > x {
// 		message = message[:x]
// 	}

// 	cli.inLine(message)

// }

// LastPath 底行刷新
func (cli *CLIOuotput) LastPath(info string, index int) {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()

	x, _ := getTerminalSize()

	message := fmt.Sprintf("[%03d] - ", index)
	message = fmt.Sprintf("%s%s", message, info)

	if len(message) > x {
		message = message[:x]
	}

	cli.inLine(message)

}

// StatusReport 状态信息打印，一般作为结果输出
func (cli *CLIOuotput) StatusReport(info string) {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()

	cli.newLine(info)
}

// inLine 本行输出打印
func (cli *CLIOuotput) inLine(output string) {
	cli.erase()
	os.Stdout.WriteString(output)
	os.Stdout.Sync()
	cli.lastInLine = true
}

// newLine 新增一行输出打印
func (cli *CLIOuotput) newLine(output string) {
	if cli.lastInLine == true {
		cli.erase()
	}

	if runtime.GOOS == "windows" {
		os.Stdout.WriteString(output)
		os.Stdout.Sync()
		os.Stdout.WriteString("\n")
		os.Stdout.Sync()
	} else {
		os.Stdout.WriteString(output + "\n")
	}
}

func percentage(x, y int) float32 {
	return float32(x) / float32(y) * 100
}
