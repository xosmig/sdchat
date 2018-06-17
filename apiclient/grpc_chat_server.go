package apiclient

import (
	"fmt"
	"github.com/xosmig/sdchat/proto"
	"google.golang.org/grpc"
	"net"
	"time"
)

type serverStream proto.Node_RouteChatServer

type grpcChatServer struct {
	grpcServer *grpc.Server
	stream     serverStream
	server     *nodeServer
}

// NewGrpcChatServer creates an instance of ApiClient
// which is used to do all the api calls from the server side.
func NewGrpcChatServer(port uint16) (ApiClient, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	nodeServer := &nodeServer{streams: make(chan serverStream)}
	proto.RegisterNodeServer(grpcServer, nodeServer)
	go grpcServer.Serve(lis)

	return &grpcChatServer{grpcServer: grpcServer, stream: nil, server: nodeServer}, nil
}

func (client *grpcChatServer) Start() error {
	client.stream = <-client.server.streams
	return nil
}

func (client *grpcChatServer) Stop() {
	client.grpcServer.Stop()
}

func (client *grpcChatServer) SendMessage(message *proto.Message) error {
	return client.stream.Send(message)
}

func (client *grpcChatServer) ReceiveMessage() (*proto.Message, error) {
	return client.stream.Recv()
}

// nodeServer implements proto.NodeServer.
// It accepts only one connection and simply saves the message stream for others' use.
type nodeServer struct {
	streams chan serverStream
	done    bool
}

func (server *nodeServer) RouteChat(stream proto.Node_RouteChatServer) error {
	if server.done {
		return fmt.Errorf("only one connection is supported")
	}

	server.streams <- stream
	server.done = true
	time.Sleep(time.Hour)
	return nil
}
