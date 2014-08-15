// +build example-main

package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	v, err := ioutil.ReadAll(os.Stdin)
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
