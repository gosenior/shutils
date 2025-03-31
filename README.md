# shutils
ExecuteScriptSh executes a shell script using the sh interpreter with arguments and output handling.

## Function ExecuteScriptSh
### Function signature
```go
func ExecuteScriptSh(filePathWdArgs []string, monitoringMode int) (string, error)
```
The script path must be relative to the current working directory of the process (where the program is executed from).
**Parameters:**
 - filePathWdArgs: []string containing:
   - [0]: Relative path to the .sh file (from where the program is executed)
   - [1:]: Optional arguments to pass to the script (can be empty)
- MonitoringMode: int specifying output handling:
   (*The function supports three execution modes, relative to the output parameter)*
   - 0: Silent mode - Runs without capturing or displaying output
   - 1: Realtime mode - Streams output with [STDOUT]/[STDERR] prefixes
   - 2: Buffered mode - Captures the combined output (stdout + stderr) for return

**Returns:**
   (output behavior)
   - string: In mode 2: Complete combined output (stdout + stderr) as a string.
   Other modes: Empty string
   - error: Runtime error, if any. For mode 2, include the captured output in the error message

**Notes:**
   - Scripts are executed with sh (not bash/zsh) - ensures compatibility with sh
   - Realtime mode (1) prefixes each line of output with either [STDOUT] or [STDERR]
   - Buffered mode (2) returns the combined output with stdout/stderr merged
   - Relative paths are resolved from the current working directory of the process

**Usage example:**

```go
func main() {
	filePath := "/scripts/script.sh"
	arg0 := "argument-0"
	arg1 := "argument-1"
	filePathWdArgs := []string{filePath, arg0, arg1}
	output, err := ExecuteScriptSh(filePathWdArgs, 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)
}
```
