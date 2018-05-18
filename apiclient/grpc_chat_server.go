package apiclient

import (
	"github.com/xosmig/sdchat2/proto"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

type serverStream sdchat2.Node_RouteChatServer

type GrpcChatServer struct {
	grpcServer *grpc.Server
	stream     serverStream
	server     *nodeServer
}

func (client *GrpcChatServer) Start() error {
	client.stream = <-client.server.streams
	return nil
}

func (client *GrpcChatServer) Stop() {
	client.grpcServer.Stop()
}

func (client *GrpcChatServer) SendMessage(message *sdchat2.Message) error {
	return client.stream.Send(message)
}

func (client *GrpcChatServer) ReceiveMessage() (*sdchat2.Message, error) {
	return client.stream.Recv()
}

func NewGrpcChatServer(port uint16) (*GrpcChatServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	nodeServer := &nodeServer{streams: make(chan serverStream)}
	sdchat2.RegisterNodeServer(grpcServer, nodeServer)
	go grpcServer.Serve(lis)

	return &GrpcChatServer{grpcServer: grpcServer, stream: nil, server: nodeServer}, nil
}

type nodeServer struct {
	streams chan serverStream
	done    bool
}

func (server *nodeServer) RouteChat(stream sdchat2.Node_RouteChatServer) error {
	if server.done {
		return fmt.Errorf("only one connection is supported")
	}

	server.streams <- stream
	server.done = true
	time.Sleep(time.Hour)
	return nil
}
