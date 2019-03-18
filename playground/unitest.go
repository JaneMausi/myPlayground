package main

import (
	"strings"
)

// the value of test_flag is from config file, that's why it's type of string
func Test_CheckTx(test_flag string, handler func(interface{}), data interface{}) {
	if 0 == strings.Compare(test_flag, "true") {
		handler(10)
	} else {
		handler(data)
	}
}
