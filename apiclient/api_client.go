package apiclient

import (
	"github.com/xosmig/sdchat/proto"
)

// ApiClient is used to do all the api calls (i.e. all the network communication).
type ApiClient interface {
	// Start must be called before any other method
	Start() error
	// Stop must be called after the use
	Stop()
	SendMessage(*proto.Message) error
	ReceiveMessage() (*proto.Message, error)
}
