// +build example-main

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	resultSeparator := "---" + "\n"

	var buf bytes.Buffer
	var err error
	var cmd, s string

	in := bufio.NewReader(os.Stdin)

GENERATION:
	for {
		buf.Reset()

		cmd, err = in.ReadString('\n')
		if err != nil {
			break GENERATION
		}
		switch cmd {
		case "Generation\n":
			// ok
		case "Exit\n":
			break GENERATION
		default:
			err = fmt.Errorf("Unknown command %q", cmd)

			break GENERATION
		}

		s, err = in.ReadString('\n')

		for err == nil && s != resultSeparator {
			_, err = buf.WriteString(s)
			if err != nil {
				break
			}

			s, err = in.ReadString('\n')
		}

		if err != nil {
			break GENERATION
		}

		fmt.Fprintf(os.Stderr, "Got step:\n%s%s\n", buf.String(), resultSeparator)

		if strings.Contains(buf.String(), "2") {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}

	if err != nil && err != io.EOF {
		panic(err)
	}

	os.Exit(0)
}
