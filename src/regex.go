package main

import (
	"fmt"
	"regexp"
)

func find(rx string, str string) (r string) {
	temp, err := regexp.Compile(rx)
	if err != nil {
		fmt.Printf("[warning] invalid regex %q: %v\n", rx, err)
		return ""
	}
	r = temp.FindString(str)
	return
}
