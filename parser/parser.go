package parser

import (
	"os"
	"strings"
	"lmd"
)

var log lmd.Logger = lmd.Logger{File: os.Stdout, FileErr: os.Stderr}

type Option int

type Argument struct {
	Name   string
	Type   Option
	Parser func(string, string)
}

func (argument *Argument) Parse(otht string) {
	argument.Parser(argument.Name, otht)
}

func Arg(Name string, Type Option, Parser func(string, string)) *Argument {
	return &Argument{
		Name:   Name,
		Type:   Type,
		Parser: Parser,
	}
}

type ArgParse struct {
	Args     []string
	Commands []*Argument
}

func quit(errcode int) {
	os.Exit(errcode)
}

func (argparse *ArgParse) Parse() {
	var recognized bool = false

	var pass bool = false

	for value, argsraw := range os.Args {

		if value == 0 {
			continue
		}

		if pass {
			pass = false
			continue
		}

		recognized = false

		for _, args := range argparse.Commands {

			switch args.Type {

			case Option(0):

				if !(len(argsraw) < 3) {

					if argsraw[:2] == "--" {

						if argsraw[2:] == strings.ToLower(args.Name) {

							args.Parse("")

							recognized = true

						} /* else {


							log  := lmd.Logger{File: os.Stdout, FileErr: os.Stderr}

							log.Log(true, "Argument " + argsraw + " is not recognized")


						} */

					}

				}

			case Option(2):

				if !(len(argsraw) < 3) {

					if argsraw[:2] == "--" {

						if argsraw[2:] == strings.ToLower(args.Name) {

							if len(os.Args) > value+1 {

								if !(len(os.Args[value+1]) < 2) {

									if os.Args[value+1][:2] != "--" {

										recognized = true

										pass = true

										args.Parse(os.Args[value+1])

									} else {

										log.Log(true, args.Name+" needs an argument")

										quit(127)

									}

								} else {

									recognized = true

									pass = true

									args.Parse(os.Args[value+1])

								}

							} else {

								log.Log(true, args.Name+" needs an argument")

								quit(127)

							}

						}

					}

				}

			case Option(3):

				if argsraw == strings.ToLower(args.Name) {

					if len(os.Args) > value+1 {

						if !(len(os.Args[value+1]) < 2) {

							if os.Args[value+1][:2] != "--" {

								recognized = true

								pass = true

								args.Parse(os.Args[value+1])

							} else {

								log.Log(true, args.Name+" needs an argument")

								quit(127)

							}

						} else {

							recognized = true

							pass = true

							args.Parse(os.Args[value+1])

						}

					} else {

						log.Log(true, args.Name+" needs an argument")

						quit(127)

					}

				}

			}

		}

		if !recognized {
			log := lmd.Logger{File: os.Stdout, FileErr: os.Stderr}
			log.Log(true, "Argument "+argsraw+" is not recognized")
			quit(127)
		}

	}
}