package main

import (
	"fmt"
	"github.com/xosmig/sdchat/apiclient"
	"github.com/xosmig/sdchat/chatnode"
	"math/rand"
	"net"
	"os"
	"time"
)

func runChatNode(conf Params) error {
	var apiClient apiclient.ApiClient
	if conf.serverIp != "" {
		// client mode
		serverIp, err := net.ResolveIPAddr("ip", conf.serverIp)
		if err != nil {
			return fmt.Errorf("cannot resolve ip address '%s': %v", conf.serverIp, err)
		}
		grpcChatClient, err := apiclient.NewGrpcChatClient(serverIp, conf.port)
		if err != nil {
			return fmt.Errorf("error connecting to the server: %v", err)
		}
		apiClient = grpcChatClient
	} else {
		// server mode
		grpcChatServer, err := apiclient.NewGrpcChatServer(conf.port)
		if err != nil {
			return fmt.Errorf("cannot initialize server: %v", err)
		}
		apiClient = grpcChatServer
	}

	chatNode := sdchat.NewChatNode(conf.name, apiClient)
	err := chatNode.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	params, err := ParseCommandLine(os.Args[1:], os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		PrintUsage(os.Stderr)
		os.Exit(2)
	}

	err = runChatNode(params)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	bye := []string{"Tschüss!", "Bye!", "Пока!"}
	fmt.Println(bye[rnd.Intn(len(bye))])
}
