package apiclient

import (
	"github.com/xosmig/sdchat2/proto"
)

type ApiClient interface {
	Start() error
	Stop()
	SendMessage(*sdchat2.Message) error
	ReceiveMessage() (*sdchat2.Message, error)
}

