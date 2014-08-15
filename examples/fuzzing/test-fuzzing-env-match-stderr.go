// +build example-main

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	f := os.Getenv("TAVOR_FUZZ_FILE")

	if f == "" {
		panic("No TAVOR_FUZZ_FILE defined")
	}

	v, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	s := string(v)

	for _, c := range s {
		fmt.Fprintf(os.Stderr, "Got %c\n", c)
	}

	os.Exit(0)
}
