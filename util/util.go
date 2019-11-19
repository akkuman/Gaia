package util

import (
	"io/ioutil"
	"strings"
	"time"
)

// TimeOut general timeout
var TimeOut = time.Duration(5 * time.Second)

// fileToStringList generate list from a dic file
func fileToStringList(path string) (strList []string, err error) {
	byteContent, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	strContent := strings.TrimSpace(string(byteContent))
	for _, line := range strings.Split(strContent, "\n") {
		lineStrip := strings.TrimSpace(line)
		if lineStrip == "" {
			continue
		}
		strList = append(strList, lineStrip)
	}
	return
}

// GetConfigStrList generate string list via a config value
func GetConfigStrList(value string) (strList []string, err error) {
	if strings.HasPrefix(value, "[file]") {
		filePath := strings.SplitN(value, ":", 2)[1]
		strList, err = fileToStringList(filePath)
		if err != nil {
			return
		}
		return
	}
	for _, v := range strings.Split(value, ",") {
		item := strings.TrimSpace(v)
		strList = append(strList, item)
	}
	return
}
