package sdchat
//go:generate protoc --go_out=plugins=grpc:. proto/message.proto
//go:generate mockgen -destination ./proto/mock_proto/node_client.mock.go github.com/xosmig/sdchat/proto NodeClient,Node_RouteChatClient
//go:generate mockgen -destination ./apiclient/mock_apiclient/api_client.mock.go github.com/xosmig/sdchat/apiclient ApiClient
