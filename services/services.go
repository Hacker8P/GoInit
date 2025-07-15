package services

import (
	"encoding/json"
	"lmd"
	"os"
	"os/exec"
	"proc"
	"strings"
	"syscall"
	"os/user"
	"strconv"
	"time"
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
	cmd := exec.Command(strings.Fields(proc.Asserter[string](service.Command))[0], strings.Fields(proc.Asserter[string](service.Command))[1:]...)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: func() uint32 {
			v, _ := strconv.Atoi(service.User.Uid)
			return uint32(v)
		}(),
		Gid: func() uint32 {
			v, _ := strconv.Atoi(service.User.Gid)
			return uint32(v)
		}(),
		Groups: func() []uint32 {
			v, _ := strconv.Atoi(service.User.Gid)
			return []uint32{uint32(v)}
		}(),
	}

	service.Status = cmd
}

/* func (ros *ReadService) ListOfServices(services []map[string]interface{}) []*Service {
	exit_list := []*Service{}

	for _, s := range services {
		exit_list = append(exit_list, MkService(proc.Asserter[string](s["Name"]), proc.Asserter[string](s["Command"]), proc.Asserter[bool](s["Active"]), proc.Asserter[int](s["At"])))
	}

	return exit_list
} */

type Array struct {
	Name string
	Command string
	Active bool
	User string
	At int
}

func (ros *ReadService) ReadServices(services *[]Service) {
	dir, _ := os.ReadDir(ros.Directory + "services")
	for _, path := range dir {
		if !path.IsDir() {
			file := ros.Directory + "services/" + path.Name()
			fl, _ := os.ReadFile(file)
			var array Array
			json.Unmarshal(fl, &array)
			service := MkService(array.Name, array.Command, array.User, array.Active, array.At)
			*services = append(*services, *service)
		}
	}
}

func (service *Service) Run() error {
	// proc.StartProcess(proc.Asserter[string](service.Name), strings.Fields(proc.Asserter[string](service.Command))[0], strings.Fields(proc.Asserter[string](service.Command))[1:], os.Stdout, os.Stderr)

	err := service.Status.Start()

	return err
}

func SetTimeout[RTYPE any](seconds int, k func() RTYPE) chan RTYPE {
	ch := make(chan RTYPE, 1)
	
	go func(ch chan RTYPE) {
		time.Sleep(time.Second * time.Duration(seconds))
		ch <- k()
	}(ch)

	return ch
}

func (service *Service) Kill() chan error {
	service.Status.Process.Signal(os.Interrupt)

	if service.Status.Process.Signal(syscall.Signal(0)) != nil {
		return nil
	}

	err := SetTimeout(2, service.Status.Process.Kill)

	return err
}

func getuserfromname(name string) (*user.User, error) {
	usr, err := user.Lookup(name)

	if err != nil {
		usrc, err := user.Current()

		return usrc, err
	}

	return usr, nil
}

func MkServiceFF(filename string) *Service {
	file, err := os.ReadFile(filename)
	var array Array
	log.Log(err != nil, err)
	err = json.Unmarshal(file, &array)
	log.Log(err != nil, err)
	return MkService(array.Name, array.Command, array.User, array.Active, array.At)
}

type StdPIPE struct {
	Stdin *PIPE
	Stdout *PIPE
	Stderr *PIPE
}

type PIPE struct {
	w *os.File
	r *os.File
}

func MkPipe() *PIPE {
	r, w, err := os.Pipe()

	log.ErrLog(err)

	return &PIPE{
		w: w,
		r: r,
	}
}

func MkStdPIPE() *StdPIPE {
	Stdin := MkPipe()
	Stdout := MkPipe()
	Stderr := MkPipe()

	return &StdPIPE{
		Stdin: Stdin,
		Stdout: Stdout,
		Stderr: Stderr,
	}
}

func MkService(name, command, user string, active bool, at int) *Service {
	usr, err := getuserfromname(user)

	if err != nil {
		usr = nil
	}

	// pipe := MkPipe()

	log.ErrLog(err)

	pipe := MkStdPIPE()

	service := &Service{
		Name: name,
		Command: command,
		Active: active,
		User: usr,
		Pipe: pipe,
		At: at,
	}

	service.makeprocess()

	service.Status.Stdout = service.Pipe.Stdout.w
	service.Status.Stdin = service.Pipe.Stdin.r
	service.Status.Stderr = service.Pipe.Stderr.w

	return service
}

type TTY struct {
	Path string
}

func RedStdPIPE() {
	
}

func read(file string) (string, error) {
	rb, err := os.ReadFile(file)

	srb := string(rb)

	return srb, err
}

func (service *Service) Attach(tty TTY) error {
	return nil
}

func MkServicePipe() {
	filer, filew, err := os.Pipe()

	log.ErrLog(err)

	cmd := exec.Command("python3")

	cmd.Stdout = filew
	cmd.Stdin = filer
	cmd.Stderr = filew

	cmd.Start()
}

type Service struct {
	Name string
	Command string
	Status *exec.Cmd
	Active bool
	User *user.User
	Pipe *StdPIPE
	At int
}