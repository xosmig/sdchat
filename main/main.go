package main

import (
	"flag"
	"os"
	"math"
	"log"
	"io"
	"github.com/xosmig/sdchat2/apiclient"
	"github.com/xosmig/sdchat2"
	"net"
)

type config struct {
	name     string
	serverIp string
	port     uint16
}

func parseArguments(errorStream io.Writer) config {
	logger := log.New(errorStream, "", 0)

	flag.Usage = func() {
		logger.Println("Usage: sdchat [-serverip IP] -port PORT NAME")
		logger.Println("If -serverip parameter is provided, the chat will run " +
			"in client mode, which means that it will connect to the server running on IP:PORT. " +
			"Otherwise, it will start a new chat server listening on port PORT.")
	}

	serverIpStrFlag := flag.String("serverip", "", "")
	portFlag := flag.Uint("port", 0, "")

	flag.Parse()

	if *portFlag == 0 {
		logger.Fatal("Parameter -port is required\n")
	}
	if *portFlag > math.MaxUint16 {
		logger.Fatalf("Invalid port: %d\n", *portFlag)
	}

	if flag.NArg() < 1 {
		logger.Print("Name is not provided\n")
		flag.Usage()
		os.Exit(3)
	}

	if flag.NArg() > 1 {
		logger.Fatalf("Unexpected argument: '%s'\n", flag.Arg(1))
	}

	return config{name: flag.Arg(0), serverIp: *serverIpStrFlag, port: uint16(*portFlag)}
}

func runChatNode(conf config) {
	var apiClient apiclient.ApiClient
	if conf.serverIp != "" {
		// client mode
		serverIp, err := net.ResolveIPAddr("ip", conf.serverIp)
		if err != nil {
			log.Fatalf("Cannot resolve ip address '%s': %v", conf.serverIp, err)
		}
		grpcChatClient, err := apiclient.NewGrpcChatClient(serverIp, conf.port)
		if err != nil {
			log.Fatalf("Error connecting to the server: %v", err)
		}
		apiClient = grpcChatClient
	} else {
		// server mode
		grpcChatServer, err := apiclient.NewGrpcChatServer(conf.port)
		if err != nil {
			log.Fatalf("Cannot initialize grpc server: %v", err)
		}
		apiClient = grpcChatServer
	}

	chatNode := sdchat.NewChatNode(conf.name, apiClient)
	chatNode.Run()
}

func main() {
	config := parseArguments(os.Stderr)
	runChatNode(config)
}
