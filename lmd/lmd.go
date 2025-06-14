package lmd

import (
	"fmt"
	"os"
	"github.com/fatih/color"
	"time"
)

func (logger *Logger) Log(is_a_err bool, logv interface{}) {
	if logv == nil {
		return
	}
	if is_a_err {
		log := "[ "
		log += color.RGB(255, 0, 0).Sprint("ERROR")

		log += " ] "
		/* _, c := i.(error)
		if c {
			fmt.Println("NOT AN ERROR")
			fmt.Fprintln(os.Stderr, "Failed wlogv)
		} */
		fmt.Fprint(logger.FileErr, log)

		fmt.Fprintln(logger.FileErr, logv)
	} else {
		fmt.Fprint(logger.File, "[ ")
		fmt.Fprint(logger.File, color.RGB(0, 255, 0).Sprint("OK"))
		fmt.Fprint(logger.File, " ] ")
		fmt.Fprintln(logger.File, logv)
	}
}

func (logger *Logger) LogInline(is_a_err bool, logv interface{}) {
	if is_a_err {
		fmt.Fprint(os.Stderr, logv)
	} else {
		fmt.Print(logv)
	}
}

func (logger *Logger) LogTime(is_a_err bool, logv string, args []any) {
	ctime := time.Now()

	timec := fmt.Sprintf("%d/%d/%d %d:%d:%d", ctime.Day(), ctime.Month(), ctime.Year(), ctime.Hour(), ctime.Minute(), ctime.Second())

	args = append(args, timec)

	if is_a_err {
		logger.Log(true, fmt.Sprintf(logv+"at %v", args...))
	} else {
		logger.Log(false, fmt.Sprintf(logv+"at %v", args...))
	}
}

type Logger struct {
	File *os.File
	FileErr *os.File
}