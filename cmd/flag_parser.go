package main

import (
	"io"
	"flag"
	"math"
	"os"
	"fmt"
)

type Params struct {
	name     string
	serverIp string
	port     uint16
}

func PrintUsage(stderr io.Writer) {
	fmt.Fprintln(stderr, "Usage: sdchat [-serverip IP] -port PORT NAME")
	fmt.Fprintln(stderr, "If -serverip parameter is provided, the chat will run " +
		"in client mode, which means that it will connect to the server running on IP:PORT. " +
		"Otherwise, it will start a new chat server listening on port PORT.")
}

func ParseCommandLine(args []string, stderr io.Writer) (Params, error) {
	commandLine := flag.NewFlagSet("sdchat", flag.ContinueOnError)
	commandLine.SetOutput(stderr)

	commandLine.Usage = func() {}

	serverIpStrFlag := commandLine.String("serverip", "", "")
	portFlag := commandLine.Uint("port", 0, "")

	err := commandLine.Parse(args)
	if err != nil {
		return Params{}, err
	}

	if *portFlag == 0 {
		return Params{}, fmt.Errorf("parameter -port is required")
	}
	if commandLine.NArg() < 1 {
		return Params{}, fmt.Errorf("parameter NAME is required")
		commandLine.Usage()
		os.Exit(3)
	}

	if *portFlag > math.MaxUint16 {
		return Params{}, fmt.Errorf("invalid port: %d", *portFlag)
	}

	if commandLine.NArg() > 1 {
		return Params{}, fmt.Errorf("Unexpected argument: '%s'\n", commandLine.Arg(1))
	}

	return Params{name: commandLine.Arg(0), serverIp: *serverIpStrFlag, port: uint16(*portFlag)}, nil
}
