package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"Gaia/plugin"
	_ "Gaia/plugin/ftp"
	_ "Gaia/plugin/ssh"
	_ "Gaia/plugin/rdp"
)

var version = "0.1"

var banner = `
_________      _____        
__  ____/_____ ___(_)_____ _
_  / __ _  __ ` + "`" + `/_  /_  __ ` + "`" + `/
/ /_/ / / /_/ /_  / / /_/ / 
\____/  \__,_/ /_/  \__,_/  
`

var (
	ips          string
	usernameFile string
	passwordFile string
	onOptions    string
	threadNum    int
)

func init() {
	flag.StringVar(&ips, "h", "", "set `host` to blast")
	flag.StringVar(&onOptions, "s", "ftp", "select the `service` to blast")
	flag.IntVar(&threadNum, "t", 10, "set `thread` num to blast")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if len(ips) == 0 || len(onOptions) == 0 {
		flag.Usage()
		return
	}

	var onOptionList []string
	for _, v := range strings.Split(onOptions, ",") {
		v = strings.TrimSpace(v)
		onOptionList = append(onOptionList, v)
	}
	pluginConfig := plugin.Config{
		IPStr:     ips,
		Options:   onOptionList,
		ThreadNum: threadNum,
	}
	pluginConfig.Start()
}

func usage() {
	fmt.Fprintf(os.Stderr, `Gaia version: Gaia/%s
%s
Options:
`, version, banner)
	flag.PrintDefaults()
}
