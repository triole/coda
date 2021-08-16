package main

import "regexp"

func find(rx string, str string) (r string) {
	temp, _ := regexp.Compile(rx)
	r = temp.FindString(str)
	return
}
