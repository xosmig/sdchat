package apiclient

import (
	"context"
	"fmt"
	"github.com/xosmig/sdchat/proto"
	"google.golang.org/grpc"
	"net"
)

type clientStream proto.Node_RouteChatClient

type GrpcChatClient struct {
	grpcClient proto.NodeClient
	conn       *grpc.ClientConn
	stream     clientStream
}

func NewGrpcChatClient(addr *net.IPAddr, port uint16) (*GrpcChatClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcChatClient{conn: conn, grpcClient: proto.NewNodeClient(conn)}, nil
}

func (client *GrpcChatClient) Start() error {
	stream, err := client.grpcClient.RouteChat(context.Background())
	if err != nil {
		return err
	}
	client.stream = stream
	return nil
}

func (client *GrpcChatClient) Stop() {
	// errors are ignored
	client.stream.CloseSend()
	client.conn.Close()
}

func (client *GrpcChatClient) SendMessage(message *proto.Message) error {
	return client.stream.Send(message)
}

func (client *GrpcChatClient) ReceiveMessage() (*proto.Message, error) {
	return client.stream.Recv()
}
