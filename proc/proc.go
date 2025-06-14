package proc

import (
	"os/exec"

	"lmd"

	// "fmt"

	"os"

	"strconv"
)

var logger_file *os.File =  func() *os.File {
	v, _ := os.OpenFile("/tmp/log/init.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)

	return v
}()

var loggermodule lmd.Logger = lmd.Logger{File: logger_file, FileErr: logger_file}

func ForceKill(process *exec.Cmd) {
	process.Process.Kill()
}

func WaitForTerminate(process *exec.Cmd, name string, ) {
	// Wait for the process's ending

	err := process.Wait()

	if err == nil {
		loggermodule.LogTime(false, "Process %v with PID %v terminated ", []any{name, strconv.Itoa(process.Process.Pid)})
	} else {
		loggermodule.LogTime(true, "Process %v with PID %v exited ", []any{name, strconv.Itoa(process.Process.Pid)})
	}
}

func Asserter[DATA any](assert_object any) DATA {
	type_asserted, _ := assert_object.(DATA)
	return type_asserted
}

func StartProcess(name string, command string, args []string, stdout *os.File, stderr *os.File) {
	proc := exec.Command(command, args...)

	// Stdout and Stderr

	proc.Stdout = stdout
	proc.Stderr = stderr

	// Starting the process

	err := proc.Start()

	// Log the starting

	/* ctime := time.Now()

	timec := fmt.Sprintf("%d/%d/%d %d:%d:%d", ctime.Day(), ctime.Month(), ctime.Year(), ctime.Hour(), ctime.Minute(), ctime.Second())

	lmd.Log(false, fmt.Sprintf("Process %v with PID %v started at %v", command, proc.Process.Pid, timec)) */

	loggermodule.LogTime(false, "Process %v with PID %v started ", []any{command, proc.Process.Pid})

	// Checking errors

	loggermodule.Log(err != nil, err)

	// go func() {err = proc.Wait(); log(err != nil, err)}()

	go WaitForTerminate(proc, command)
}