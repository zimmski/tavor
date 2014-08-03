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

	in := bufio.NewReader(os.Stdin)

	s, err := in.ReadString('\n')

	for err == nil && s != resultSeparator {
		_, err = buf.WriteString(s)
		if err != nil {
			break
		}

		s, err = in.ReadString('\n')
	}

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Got original input:\n%s%s\n", buf.String(), resultSeparator)

	fmt.Println("OK")

	for {
		buf.Reset()

		s, err = in.ReadString('\n')

		for err == nil && s != resultSeparator {
			_, err = buf.WriteString(s)
			if err != nil {
				break
			}

			s, err = in.ReadString('\n')
		}

		if err != nil {
			break
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
