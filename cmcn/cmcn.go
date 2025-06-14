package cmcn

import (
	"golang.org/x/sys/unix"

	"lmd"

	"os"

	"strings"

	"proc"

	"smng"
)

func createFifoIfNotExist() {
	unix.Mkfifo("/tmp/goinit_fifo", 0644)
}

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
			lmd.LogInline(false, msg_i)
	}
}

func wait_for_message() {
	for {
		m, _ := os.OpenFile("/tmp/goinit_fifo", os.O_RDONLY, os.ModeNamedPipe)

		l := make([]byte, 128)

		c, _ := m.Read(l)

		m.Close()

		elaborate(string(l[:c]))
	}
}