package services

import (
	"encoding/json"
	"errors"
	"io"
	"lmd"
	"os"
	"os/exec"
	"os/user"
	"proc"
	"strconv"
	"strings"
	"syscall"
	"time"
	"github.com/creack/pty"
	"golang.org/x/term"
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

func (service *Service) RunPTPipe() error {
	pipe, err := pty.Start(service.Status)

	log.ErrLog(err)

	service.PTPipe = pipe

	return err
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

	// service.Status.Stdout = service.Pipe.Stdout.w
	// service.Status.Stdin = service.Pipe.Stdin.r
	// service.Status.Stderr = service.Pipe.Stderr.w

	return service
}

type TTY struct {
	Path string
	File *os.File
	O *term.State
}

/* func read(file string) (string, error) {
	rb, err := os.ReadFile(file)

	srb := string(rb)

	return srb, err
} */

func asyncStdinPipe(pipe *os.File, tty *TTY) {
	_, err := io.Copy(pipe, tty.File)

	log.Log(false, "Finished!!")

	log.ErrLog(err)
}

func asyncStdoutPipe(pipe *os.File, tty *TTY) {
	_, err := io.Copy(tty.File, pipe)

	log.ErrLog(err)
}

func asyncStderrPipe(pipe *os.File, tty *TTY) {
	_, err := io.Copy(tty.File, pipe)

	log.ErrLog(err)
}

func (service *Service) Attach(tty *TTY) error {
	if service.PTPipe == nil {
		return errors.New("cannot attach to a tty a process not started")
	}

	o, err := term.MakeRaw(int(tty.File.Fd()))

	tty.O = o

	if err != nil {
		return err
	}

	go asyncStdinPipe(service.PTPipe, tty)
	go asyncStdoutPipe(service.PTPipe, tty)
	// go asyncStderrPipe(service.PTPipe, tty)

	log.LogTime(false, "Attaching " + service.Name + " ", []any{})

	return nil
}

func (service *Service) Detach(tty *TTY) error {
	err := term.Restore(int(tty.File.Fd()), tty.O)

	return err
}

type ENVVAR struct {
	Key string
	Value string
}

func GetEnv(unixuser *user.User) ([]ENVVAR, error) {
	command := exec.Command("env", "-i", "bash", "--login", "-c", "env")

	command.SysProcAttr = &syscall.SysProcAttr{}

	command.SysProcAttr.Credential = &syscall.Credential{
		Uid: func() uint32 {
			v, _ := strconv.Atoi(unixuser.Uid)
			return uint32(v)
		}(),
		Gid: func() uint32 {
			v, _ := strconv.Atoi(unixuser.Gid)
			return uint32(v)
		}(),
		Groups: func() []uint32 {
			v, _ := strconv.Atoi(unixuser.Gid)
			return []uint32{uint32(v)}
		}(),
	}

	commout, err := command.Output()

	commoutstr := string(commout)

	lcs := strings.Split(commoutstr, "\n")

	var result []ENVVAR = []ENVVAR{}

	for n, h := range lcs {
		if strings.Contains(h, "=") && !strings.Contains(h, "tty") {
			k := strings.Split(h, "=")

			result[n] = ENVVAR{
				k[0],
				k[1],
			}
		}
	}

	return result, err
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
	PTPipe *os.File
	At int
}