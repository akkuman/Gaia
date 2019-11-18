package clioutput

import (
	"os"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type (
	// DWORD 双字
	DWORD int32
)

// FillConsoleOutputCharacterProc 从指定的坐标开始，将指定的次数的字符写入控制台屏幕缓冲区。
var FillConsoleOutputCharacterProc *windows.Proc

func initCLIOuotput() (clioutpt CLIOuotput) {
	mutex := new(sync.Mutex)

	kernel32DLL := windows.MustLoadDLL("kernel32.dll")
	defer kernel32DLL.Release()
	FillConsoleOutputCharacterProc = kernel32DLL.MustFindProc("FillConsoleOutputCharacterW")

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
	csbi := mustGetCsbi()
	// \b 代表退格
	// 把光标退格移到行首
	line := strings.Repeat("\b", int(csbi.CursorPosition.X))
	os.Stdout.WriteString(line)
	width := csbi.CursorPosition.X
	csbi.CursorPosition.X = 0
	out := 0
	// 把光标定位到控制台底部
	// 参考链接 https://github.com/erikh/termproxy/blob/master/dockerterm/winconsole/console_windows.go#L376
	if err := getError(FillConsoleOutputCharacterProc.Call(uintptr(windows.Stdout), ' ', uintptr(width), marshal(csbi.CursorPosition), uintptr(unsafe.Pointer(&out)))); err != nil {
		panic(err)
	}
	os.Stdout.WriteString(line)
	// 刷新控制台缓冲区
	os.Stdout.Sync()
}

// 获取控制台宽度高度
func getTerminalSize() (columns, rows int) {
	csbi := mustGetCsbi()

	columns = int(csbi.Window.Right - csbi.Window.Left)
	rows = int(csbi.Window.Bottom - csbi.Window.Top)
	return
}

// mustGetCsbi 获取csbi(ConsoleScreenBufferInfo)，出错就panic
func mustGetCsbi() windows.ConsoleScreenBufferInfo {
	hConsole := windows.Stdout
	csbi := new(windows.ConsoleScreenBufferInfo)
	err := windows.GetConsoleScreenBufferInfo(hConsole, csbi)
	if err != nil {
		panic(err)
	}
	return *csbi
}

// 参考链接：
// https://docs.microsoft.com/en-us/windows/console/fillconsoleoutputcharacter
// https://github.com/erikh/termproxy/blob/master/dockerterm/winconsole/console_windows.go#L1044
func marshal(c windows.Coord) uintptr {
	// 指向 csbi.CursorPosition(windows.Coord) 结构体的指针
	ptrCoord := (*DWORD)(unsafe.Pointer(&c))
	coord := *ptrCoord
	return uintptr(coord)
}

// 检查函数是否发生错误
// 参考链接
// https://github.com/erikh/termproxy/blob/master/dockerterm/winconsole/console_windows.go#L297
func getError(r1, r2 uintptr, lastErr error) error {
	// If the function fails, the return value is zero.
	if r1 == 0 {
		if lastErr != nil {
			return lastErr
		}
		return syscall.EINVAL
	}
	return nil
}
