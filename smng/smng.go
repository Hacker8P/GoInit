package smng

import (
	"encoding/json"
	"lmd"
	"os"
	"os/exec"
	"proc"
	"strings"
)

/* func AssertServices(services *[]map[string]interface{}) {
	for _, service := range *services {
		for _, types := range service {
			types = proc.Asserter[]()
		}
	}
} */

var logger_file *os.File =  func() *os.File {
	v, _ := os.OpenFile("/tmp/log/init.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)

	return v
}()


var log lmd.Logger = lmd.Logger{File: logger_file, FileErr: logger_file}

type ReadService struct {
	Directory string
}

func (service *Service) makeprocess() {
	service.Status = exec.Command(strings.Fields(proc.Asserter[string](service.Command))[0], strings.Fields(proc.Asserter[string](service.Command))[1:]...)
}

/* func (ros *ReadService) ListOfServices(services []map[string]interface{}) []*Service {
	exit_list := []*Service{}

	for _, s := range services {
		exit_list = append(exit_list, MkService(proc.Asserter[string](s["Name"]), proc.Asserter[string](s["Command"]), proc.Asserter[bool](s["Active"]), proc.Asserter[int](s["At"])))
	}

	return exit_list
} */

func (ros *ReadService) ReadServices(services *[]Service) {
	dir, _ := os.ReadDir(ros.Directory + "services")
	for _, path := range dir {
		if !path.IsDir() {
			file := ros.Directory + "services/" + path.Name()
			fl, _ := os.ReadFile(file)
			type Array struct {
				Name string
				Command string
				Active bool
				At int
			}
			var array Array
			json.Unmarshal(fl, &array)
			service := MkService(array.Name, array.Command, array.Active, array.At)
			*services = append(*services, *service)
		}
	}
}

func (service *Service) Run() {
	// proc.StartProcess(proc.Asserter[string](service.Name), strings.Fields(proc.Asserter[string](service.Command))[0], strings.Fields(proc.Asserter[string](service.Command))[1:], os.Stdout, os.Stderr)
	service.Status.Start()
}

func (service *Service) Kill() {
	service.Status.Process.Kill()
}

func MkServiceFF(filename string) *Service {
	file, err := os.ReadFile(filename)
	type Array struct {
		Name string
		Command string
		Active bool
		At int
	}
	var array Array
	log.Log(err != nil, err)
	err = json.Unmarshal(file, &array)
	log.Log(err != nil, err)
	return MkService(array.Name, array.Command, array.Active, array.At)
}

func MkService(name, command string, active bool, at int) *Service {
	service := &Service{
		Name: name,
		Command: command,
		Active: active,
		At: at,
	}

	service.makeprocess()

	return service
}

type Service struct {
	Name string
	Command string
	Status *exec.Cmd
	Active bool
	At int
}