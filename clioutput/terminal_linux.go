package clioutput

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

// ESC esc按键 ASNI/VT100
var ESC = "\x1b"

func initCLIOuotput() (clioutpt CLIOuotput) {
	mutex := new(sync.Mutex)

	clioutpt = CLIOuotput{
		lastLength: 0,
		lastOutput: "",
		lastInLine: false,
		mutex:      mutex,
	}
	return
}

// 参考链接 https://github.com/maurosoria/dirsearch/blob/master//lib/output/CLIOutput.py#L52
func (cli *CLIOuotput) erase() {
	// Erase Start of Line
	os.Stdout.WriteString(fmt.Sprintf("%s%s", ESC, "[1K"))
	// Clear a tab at the current column
	os.Stdout.WriteString(fmt.Sprintf("%s%s", ESC, "[0G"))
}

// 获取unix控制台宽度高度
// 参考链接 https://github.com/golang/crypto/blob/master/ssh/terminal/util.go#L79
func getTerminalSize() (columns, rows int) {
	fd := int(os.Stdin.Fd())
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}
	return int(ws.Col), int(ws.Row)
}
