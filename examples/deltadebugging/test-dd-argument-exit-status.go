// +build example-main

package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		panic("No TAVOR_DD_FILE argument")
	}

	v, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	s := string(v)

	if strings.Contains(s, "2") {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
