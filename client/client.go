package main

import (
	"encoding/json"

	"strings"

	"fmt"

	"lmd"

	"os"

	"parser"
)

var log lmd.Logger = lmd.Logger{File: os.Stdout, FileErr: os.Stderr}

var help string = `GOINIT

NAME
       GoInit - system and service manager

SYNOPSIS

	goinit [OPTIONS...]

	goinit [COMMANDS]

COMMANDS

	start [SERVICE]

	stop [SERVICE]

	restart [SERVICE]

	enable [SERVICE]

	disable [SERVICE]
`

// var help string = "GOINIT\n\nNAME\n       GoInit - system and service manager\n\nSYNOPSIS\n\n\t/usr/lib/systemd/systemd [OPTIONS...]\n\n\tgoinit [COMMANDS]\n       \nCOMMANDS\n\n\tstart [SERVICE]\n\t\n\tstop [SERVICE]\n\t\n\trestart [SERVICE]\n\t\n\tenable [SERVICE]\n\t\n\tdisable [SERVICE]\n"

var channel chan error = make(chan error)

var jsonstruct JSON = JSON{File: FIFO{"/tmp/goinit_fifo"}, Data: nil}

type JsonData struct {
	Command string
}

type FIFO struct {
	FilePath string
}

func (fifo *FIFO) Write(message []byte, channel chan error) {
	r, err := os.OpenFile(fifo.FilePath, os.O_WRONLY, os.ModeNamedPipe)
	r.Write([]byte(message))
	defer r.Close()

	channel <- err
}

func (fifo *FIFO) Name() string {
	return fifo.FilePath
}

type JSON struct {
	File FIFO
	Data *json.Marshaler
}

func (jstruct *JSON) Write(message JsonData, channel chan error) error {
	jsondatamarshal, err := json.Marshal(message)

	log.Log(err != nil, err)

	go jstruct.File.Write(jsondatamarshal, channel)

	err = <-channel

	return err
}

func quit(errcode int) {
	os.Exit(errcode)
}

func Parser(name string, otht string) {
	switch strings.ToLower(name) {
	case "help":
		fmt.Println(help)
		quit(0)
	case "version":
		fmt.Println("VERSION")
		quit(0)
	default:
		go jsonstruct.Write(JsonData{Command: strings.ToLower(name) + " " + otht}, channel)

		err := <-channel

		log.Log(err != nil, err)
	}
}

func main() {
	Arg := parser.Arg

	type ArgParse = parser.ArgParse

	type Argument = parser.Argument

	type Option = parser.Option

	argparse := ArgParse{Args: os.Args, Commands: []*Argument{
		Arg("Help", Option(0), Parser),
		Arg("Version", Option(0), Parser),
		Arg("Start", Option(3), Parser),
		Arg("Restart", Option(3), Parser),
		Arg("Enable", Option(3), Parser),
		Arg("Disable", Option(3), Parser),
		Arg("Stop", Option(3), Parser),
	}}

	argparse.Parse()

	/* jstruct_example := JSON{File: FIFO{"/tmp/goinit_fifo"}, Data: nil}

	channel := make(chan error) */

	/* go jstruct_example.Write(JsonData{Command: "Start Ciao"}, channel)

	err := <-channel

	fmt.Println(err) */
}
