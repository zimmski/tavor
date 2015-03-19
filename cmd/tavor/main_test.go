package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	f, err := ioutil.TempFile("", "tavor-main-test")
	assert.Nil(t, err)

	_, err = f.WriteString("START = 1 | 2 | 3\n")
	assert.Nil(t, err)

	err = f.Close()
	assert.Nil(t, err)

	defer func() {
		err := os.Remove(f.Name())
		assert.Nil(t, err)
	}()

	exitCode, out := execMain(t, []string{"--format-file", f.Name(), "fuzz", "--strategy", "AllPermutations"})

	assert.Equal(t, exitCodeOk, exitCode)
	assert.Contains(t, out, "1\n2\n3\n")
}

func TestMainCommandListingOptions(t *testing.T) {

	exitCode, out := execMain(t, []string{"fuzz", "--list-exec-argument-types"})

	assert.Equal(t, exitCodeOk, exitCode)
	assert.Contains(t, out, strings.Join(execArgumentTypes, "\n"))
}

func execMain(t *testing.T, args []string) (exitCodeType, string) {
	saveStderr := os.Stderr
	saveStdout := os.Stdout
	saveCwd, err := os.Getwd()
	assert.Nil(t, err)

	r, w, err := os.Pipe()
	assert.Nil(t, err)

	os.Stderr = w
	os.Stdout = w

	bufChannel := make(chan string)

	go func() {
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, r)
		assert.Nil(t, err)
		assert.Nil(t, r.Close())

		bufChannel <- buf.String()
	}()

	exitCode := mainCmd(args)

	assert.Nil(t, w.Close())

	os.Stderr = saveStderr
	os.Stdout = saveStdout
	assert.Nil(t, os.Chdir(saveCwd))

	out := <-bufChannel

	return exitCode, out
}
