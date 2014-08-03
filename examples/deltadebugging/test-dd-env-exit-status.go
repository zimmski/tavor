// +build example-main

package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	f := os.Getenv("TAVOR_DD_FILE")

	if f == "" {
		panic("No TAVOR_DD_FILE defined")
	}

	v, err := ioutil.ReadFile(f)
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
