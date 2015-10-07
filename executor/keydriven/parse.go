package keydriven

import (
	"io/ioutil"
	"strings"
)

// ReadKeyDrivenFile reads in a key driven file.
// A key driven file has the following format: Each lines holds one command. A command is defined by its key and parameters which are separated by a tabulator. The line ends with a new line character.
func ReadKeyDrivenFile(file string) ([]Command, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cmds []Command

	for _, l := range strings.Split(string(input), "\n") {
		lc := strings.Split(l, "\t")

		for i := 0; i < len(lc); i++ {
			lc[i] = strings.Trim(lc[i], "\r ")
		}

		if len(lc[0]) != 0 {
			cmds = append(cmds, Command{
				Key:        lc[0],
				Parameters: lc[1:],
			})
		}
	}

	return cmds, nil
}
