package sdchat

import (
	"bufio"
	"context"
	"fmt"
	"github.com/xosmig/sdchat/apiclient"
	"github.com/xosmig/sdchat/proto"
	"github.com/xosmig/sdchat/util"
	"io"
	"log"
	"os"
	"time"
)

// ChatNode provides simple text user interface for sdchat.
type ChatNode struct {
	reader     *bufio.Reader
	clientName string
	apiClient  apiclient.ApiClient
	stdout     io.Writer
}

// NewChatNode creates a chat node with the given username and ApiClient
func NewChatNode(username string, apiClient apiclient.ApiClient) ChatNode {
	return ChatNode{bufio.NewReader(os.Stdin), username, apiClient, os.Stdout}
}

// sendMessage constructs and sends a message with the given text.
func (node *ChatNode) sendMessage(text string) {
	message := &proto.Message{
		Name:      node.clientName,
		Timestamp: time.Now().Unix(),
		Text:      text,
	}
	err := node.apiClient.SendMessage(message)
	if err != nil {
		node.printError(fmt.Sprintf("Oops: error sending message: %v", err))
		return
	}
	node.printMessage(message)
}

// printf prints the text to the node's output channel
func (node *ChatNode) printf(format string, a ...interface{}) {
	fmt.Fprintf(node.stdout, format, a...)
}

// printf prints the line to the node's output channel
func (node *ChatNode) println(line string) {
	node.printf("%s\n", line)
}

// printMessage prints the given message to the node's output channel
func (node *ChatNode) printMessage(message *proto.Message) {
	timeStr := time.Unix(message.Timestamp, 0).Format("02.01.2006 15:04")
	node.printf("\n[%s] %s: %s\n", timeStr, message.Name, message.Text)
}

// printError prints the given error message to the node's output channel
func (node *ChatNode) printError(errorMsg string) {
	node.printf("\nError: %s\n", errorMsg)
}

// Run starts the chat node.
// It blocks until one of the clients decides to close the connection, or an error happens.
// Use RunWithContext if you want to be able to interrupt the execution.
func (node *ChatNode) Run() error {
	return node.RunWithContext(context.Background())
}

// RunWithContext is the same as Run, but can be interrupted by cancelling the context.
func (node *ChatNode) RunWithContext(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Println("Connecting...")
	err := node.apiClient.Start()
	if err != nil {
		node.printError(fmt.Sprintf("Connection error: %v", err))
		return err
	}
	defer node.apiClient.Stop()
	log.Println("Connected.")
	log.Println("Type \"m[enter]\" to write a message or \"q[enter]\" to exit")

	lines := make(chan string)
	go func() {
		defer close(lines)
		for {
			bytes, isPrefix, err := node.reader.ReadLine()
			if isPrefix {
				node.printError("Line is too long")
				err = util.DiscardLineFromReader(node.reader)
			}
			if err == io.EOF {
				cancel()
				return
			}
			if err != nil {
				node.printError(fmt.Sprintf("Oops: error reading your input: %v", err))
				continue
			}
			select {
			case lines <- string(bytes):
			case <-ctx.Done():
				return
			}
		}
	}()

	messages := make(chan *proto.Message)
	go func() {
		for {
			message, err := node.apiClient.ReceiveMessage()
			if err != nil {
				node.println("Connection is closed")
				cancel()
				return
			}
			select {
			case messages <- message:
			case <-ctx.Done():
				return
			}
		}
	}()

selectLoop:
	for {
		log.Println("Debug: waiting for an event...")
		select {
		case <-ctx.Done():
			log.Println("Debug: exiting")
			break selectLoop
		case line := <-lines:
			switch line {
			case "m":
				log.Println("Debug: reading message from the user...")
				node.printf("Enter message: ")
				select {
				case <-ctx.Done():
				case text := <-lines:
					log.Println("Debug: sending message...")
					node.sendMessage(text)
				}
			case "q":
				log.Println("Debug: exit requested")
				cancel()
			default:
				node.printError(fmt.Sprintf("Unknown command: '%s'", line))
			}
		case message := <-messages:
			node.printMessage(message)
		}
	}

	return nil
}
