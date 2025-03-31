package shutils

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteScriptShOk(t *testing.T) {
	output, err := ExecuteScriptSh([]string{"scripts/script.sh", "arg1", "arg2"}, 2)
	assert.Nil(t, err)
	expected := "Received arguments: arg1 arg2\nFirst argument: arg1\nSecond argument: arg2\n"
	assert.Equal(t, expected, output)
}

func TestExecuteScriptShError(t *testing.T) {
	_, err := ExecuteScriptSh([]string{}, 0)
	assert.Contains(t, err.Error(), "filePathWdArgs cannot be empty")
}

func TestMonitoringMode0DiscardOutput(t *testing.T) {
	filePath := "/scripts/script.sh"
	arg0 := "argument-0"
	arg1 := "argument-1"
	filePathWdArgs := []string{filePath, arg0, arg1}
	output, err := ExecuteScriptSh(filePathWdArgs, 0)
	assert.Nil(t, err)
	assert.Equal(t, "", output)
}

func TestMonitoringMode1RealTimeOutput(t *testing.T) {
	var buf bytes.Buffer
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	done := make(chan struct{})
	go func() {
		_, _ = buf.ReadFrom(r)
		close(done)
	}()

	filePath := "/scripts/script.sh"
	arg0 := "argument-0"
	arg1 := "argument-1"
	filePathWdArgs := []string{filePath, arg0, arg1}

	output, err := ExecuteScriptSh(filePathWdArgs, 1)

	_ = w.Close()
	os.Stdout = origStdout
	<-done 

	assert.Nil(t, err)
	assert.Equal(t, "", output)
	assert.NotEmpty(t, buf.String(), "I expected a real-time exit, but there was no exit")
}

func TestMonitoringMode2ReturnOutput(t *testing.T) {
	filePath := "/scripts/script.sh"
	arg0 := "argument-0"
	arg1 := "argument-1"
	filePathWdArgs := []string{filePath, arg0, arg1}
	output, err := ExecuteScriptSh(filePathWdArgs, 2)
	assert.Nil(t, err)
	expected := "Received arguments: argument-0 argument-1\nFirst argument: argument-0\nSecond argument: argument-1\n"
	assert.Equal(t, expected, output)
}
