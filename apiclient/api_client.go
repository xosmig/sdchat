package apiclient

import (
	"github.com/xosmig/sdchat/proto"
)

type ApiClient interface {
	Start() error
	Stop()
	SendMessage(*proto.Message) error
	ReceiveMessage() (*proto.Message, error)
}
