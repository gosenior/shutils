package shutils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// ExecuteScriptSh executes a shell script using the sh interpreter with arguments and output handling.
//
// The function supports three execution modes: silent, real-time output, and buffered output capture.
// The script path must be relative to the current working directory of the process (where the program is executed from). //
// Parameters:
// - filePathWdArgs: []string containing:
// - [0]: Relative path to the .sh file (relative to the root directory of the mod file, i.e. the root path of the project)
// - [1:]: Optional arguments to pass to the script (can be empty)
// - monitoringMode: int specifying output handling:
// - 0: Silent mode - Runs without capturing or displaying output
// - 1: Realtime mode - Streams output with [STDOUT]/[STDERR] prefixes
// - 2: Buffered mode - Captures the combined output (stdout + stderr) for return
//
// Returns:
// - string: In mode 2: Complete combined output (stdout + stderr) as a string.
// Other modes: Empty string
// - error: Runtime error, if any. For mode 2, include the captured output in the error message
//
// Notes:
// - Scripts are executed with sh (not bash/zsh) - ensures compatibility with sh
// - Realtime mode (1) prefixes each line of output with either [STDOUT] or [STDERR]
// - Buffered mode (2) returns the combined output with stdout/stderr merged
// - Relative paths are resolved from the current working directory of the process
func ExecuteScriptSh(filePathWdArgs []string, monitoringMode int) (string, error) {
	// Basic parameter validation
	if len(filePathWdArgs) == 0 {
		return "", fmt.Errorf("filePathWdArgs cannot be empty")
	}

	// Get project's current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	// Build absolute script path
	scriptPath := filepath.Join(currentDir, filePathWdArgs[0])

	// Prepare command arguments
	args := append([]string{scriptPath}, filePathWdArgs[1:]...)

	// Create command
	cmd := exec.Command("sh", args...)

	switch monitoringMode {
	case 0:
		// Silent execution (discard output)
		err := cmd.Run()
		return "", err

	case 1:
		// Real-time output
		err := showOutputInRealTime(cmd, scriptPath)
		return "", err

	case 2:
		// Buffered output (capture and return output)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("error executing script %s: %v\nOutput:\n%s", scriptPath, err, output)
		}
		return string(output), nil

	default:
		return "", fmt.Errorf("invalid monitoringMode value: %d", monitoringMode)
	}
}

// showOutputInRealTime streams script output to stdout in real-time.
func showOutputInRealTime(cmd *exec.Cmd, scriptPath string) error {
	// Get pipes for real-time output
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting stderr pipe: %v", err)
	}

	// Start the command before scanning the output
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting script %s: %v", scriptPath, err)
	}

	// Use goroutines to read stdout and stderr in real-time
	go streamOutput(stdoutPipe, "STDOUT")
	go streamOutput(stderrPipe, "STDERR")

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error executing script %s: %v", scriptPath, err)
	}

	return nil
}

// streamOutput reads from an output pipe and prints to the console.
func streamOutput(pipe io.Reader, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}
