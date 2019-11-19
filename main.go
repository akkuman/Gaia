package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"Gaia/plugin"
	_ "Gaia/plugin/ftp"
	_ "Gaia/plugin/smb"
	_ "Gaia/plugin/ssh"
	"Gaia/util"
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
	ips         string
	userListStr string
	passListStr string
	onOptions   string
	threadNum   int
)

func init() {
	flag.StringVar(&ips, "h", "", "set `host` to blast")
	flag.StringVar(&userListStr, "u", "", "set `username` list, such as: -u admin,root OR -u [file]:username.txt")
	flag.StringVar(&passListStr, "p", "", "set `password` list, such as: -p admin,toor OR -u [file]:password.txt")
	flag.StringVar(&onOptions, "s", "ftp", "select the `service` to blast, options: ftp,smb,ssh")
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

	pluginConfig := initConfig()
	pluginConfig.Start()
}

func initConfig() (pluginConfig plugin.Config) {
	var onOptionList []string
	var userList []string
	var passList []string
	for _, v := range strings.Split(onOptions, ",") {
		onOptionList = append(onOptionList, strings.TrimSpace(v))
	}
	if strings.HasPrefix(userListStr, "[file]:") {
		userList, _ = util.GetConfigStrList(userListStr)
	} else {
		for _, v := range strings.Split(userListStr, ",") {
			userList = append(userList, strings.TrimSpace(v))
		}
	}
	if strings.HasPrefix(passListStr, "[file]:") {
		passList, _ = util.GetConfigStrList(passListStr)
	} else {
		for _, v := range strings.Split(passListStr, ",") {
			passList = append(passList, strings.TrimSpace(v))
		}
	}
	pluginConfig = plugin.Config{
		IPStr:     ips,
		Options:   onOptionList,
		ThreadNum: threadNum,
		UserList:  userList,
		PassList:  passList,
	}
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, `Gaia version: Gaia/%s
%s
Options:
`, version, banner)
	flag.PrintDefaults()
}
