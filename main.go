package main

import (
	"fmt"
	"os"
)

var command = map[string]string{
	"help": " -h, --help",
}

// HELP 帮助信息
var HELP = fmt.Sprintf(`
Usage: 

%-18s display this help and exit
`, command["help"])

func main() {
	switch len(os.Args) {
	case 1:
		help()
	case 2:
		switch os.Args[1] {
		case "-h":
			fallthrough
		case "--help":
			help()
		}
	}
}

func help() {
	fmt.Print(HELP)
}
