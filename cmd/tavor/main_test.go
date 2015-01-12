package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
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
		r.Close()
		assert.Nil(t, err)

		bufChannel <- buf.String()
	}()

	exitCode := mainCmd([]string{"--format-file", f.Name(), "fuzz", "--strategy", "AllPermutations"})

	w.Close()

	os.Stderr = saveStderr
	os.Stdout = saveStdout
	os.Chdir(saveCwd)

	out := <-bufChannel

	assert.Equal(t, exitCodeOk, exitCode)
	assert.Contains(t, out, "1\n2\n3\n")
}
