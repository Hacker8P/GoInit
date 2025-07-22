module smng

go 1.24.3

replace lmd => ../lmd

replace proc => ../proc/

require (
	github.com/creack/pty v1.1.24
	golang.org/x/term v0.33.0
	lmd v0.0.0-00010101000000-000000000000
	proc v0.0.0-00010101000000-000000000000
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.34.0 // indirect
)
