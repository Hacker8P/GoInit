package main

import (
	// "fmt"
	"os"
	// "os/exec"
	// "strconv"
	"strings"
	// "time"

	"services"

	"lmd"

	"proc"

	"golang.org/x/sys/unix"

	// "syscall"
)

func createFifoIfNotExist() {
	unix.Mkfifo("/tmp/goinit_fifo", 0644)
}

var services_files []services.Service

var logger_file *os.File =  func() *os.File {
	v, _ := os.OpenFile("/tmp/log/init.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)

	return v
}()

var logger lmd.Logger = lmd.Logger{File: logger_file, FileErr: logger_file}

func elaborate(msg_i string) {
	// log_inline(false, msg_i)
	msg := strings.Fields(msg_i)
	if len(msg) <= 0 {
		createFifoIfNotExist()
		return
	}
	switch msg[0] {
	case "poweroff":
		os.Exit(0)
	case "run":
		var args_0 []string
		if len(msg) > 3 {
			args_0 = strings.Split(msg[3], ",")
		} else {
			args_0 = []string{}
		}
		proc.StartProcess(msg[1], msg[2], args_0, os.Stdout, os.Stderr)
	case "reload":
		ServiceElaborate()
	default:
		logger.LogInline(false, msg_i)
	}
}

type JSON struct {
	Command string
}

func wait_for_message() {
	for {
		m, _ := os.OpenFile("/tmp/goinit_fifo", os.O_RDONLY, os.ModeNamedPipe)

		/* var json_output JSON

		json.Unmarshal() */

		l := make([]byte, 128)

		c, _ := m.Read(l)

		m.Close()

		elaborate(string(l[:c]))
	}
}

func ServiceElaborate() {
	ros := services.ReadService{Directory: "/etc/goinit/"}
	ros.ReadServices(&services_files)
	for _, service_single := range services_files {
		service_single.Status.Stdout = os.Stdout
		service_single.Status.Stderr = os.Stderr
		service_single.Status.Stdin = os.Stdin
		service_single.Run()
	}
}

func RunBg() {
	x, _ := os.OpenFile("/tmp/init_log", os.O_RDWR, os.ModeAppend)
	proc.StartProcess("Init", "./init", []string{"--bg"}, x, x)
}

func main() {
	/* if len(os.Args) > 1 {
		if os.Args[1] != "--bg" {
			RunBg()
			os.Exit(0)
		}
	} else {
		RunBg()
		os.Exit(0)
	} */
	// proc.StartProcess("Bash", "sleep", []string{"5"}, os.Stdout, os.Stderr)

	services.MkService("Bash", "sleep 5", "pietro", true, 0).Run()

	// communicate("Ciao")

	/* for a, _ := range services_files {
		x := services.Service{a["Name"], }
	} */

	/* if assertion, may := services_files[0].(map[string]interface{}); may {
		fmt.Println(assertion)
		fmt.Println(len(services_files))
	} */

	/* fmt.Println(services_files[0]["Argument"])

	fmt.Println(services_files_Service[0].Command) */

	ServiceElaborate()

	logger.Log(false, "Example")

	wait_for_message()
}