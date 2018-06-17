package apiclient

import (
	"context"
	"fmt"
	"github.com/xosmig/sdchat/proto"
	"google.golang.org/grpc"
	"net"
)

type clientStream proto.Node_RouteChatClient

type grpcChatClient struct {
	grpcClient proto.NodeClient
	conn       *grpc.ClientConn
	stream     clientStream
}

// NewGrpcChatClient creates an instance of ApiClient
// which is used to do all the api calls from the client side.
func NewGrpcChatClient(addr *net.IPAddr, port uint16) (ApiClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &grpcChatClient{conn: conn, grpcClient: proto.NewNodeClient(conn)}, nil
}

func (client *grpcChatClient) Start() error {
	stream, err := client.grpcClient.RouteChat(context.Background())
	if err != nil {
		return err
	}
	client.stream = stream
	return nil
}

func (client *grpcChatClient) Stop() {
	// errors are ignored
	client.stream.CloseSend()
	client.conn.Close()
}

func (client *grpcChatClient) SendMessage(message *proto.Message) error {
	return client.stream.Send(message)
}

func (client *grpcChatClient) ReceiveMessage() (*proto.Message, error) {
	return client.stream.Recv()
}
