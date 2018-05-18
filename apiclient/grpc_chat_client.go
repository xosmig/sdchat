package apiclient

import (
	"net"
	"google.golang.org/grpc"
	"fmt"
	"github.com/xosmig/sdchat2/proto"
	"context"
)

type clientStream sdchat2.Node_RouteChatClient

type GrpcChatClient struct {
	grpcClient sdchat2.NodeClient
	conn       *grpc.ClientConn
	stream     clientStream
}

func NewGrpcChatClient(addr *net.IPAddr, port uint16) (*GrpcChatClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcChatClient{conn: conn, grpcClient: sdchat2.NewNodeClient(conn)}, nil
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

func (client *GrpcChatClient) SendMessage(message *sdchat2.Message) error {
	return client.stream.Send(message)
}

func (client *GrpcChatClient) ReceiveMessage() (*sdchat2.Message, error) {
	return client.stream.Recv()
}
