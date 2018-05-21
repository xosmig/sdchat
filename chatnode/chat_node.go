package sdchat

import (
	"bufio"
	"fmt"
	"time"
	"github.com/xosmig/sdchat/apiclient"
	"os"
	"log"
	"io"
	"github.com/xosmig/sdchat/proto"
	"context"
)

type ChatNode struct {
	reader    *bufio.Reader
	name      string
	apiClient apiclient.ApiClient
	stdout    io.Writer
}

func NewChatNode(name string, apiClient apiclient.ApiClient) ChatNode {
	return ChatNode{bufio.NewReader(os.Stdin), name, apiClient, os.Stdout}
}

func (node *ChatNode) sendMessage(text string) {
	message := &proto.Message{
		Name:      node.name,
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

func (node *ChatNode) printf(format string, a ...interface{}) {
	fmt.Fprintf(node.stdout, format, a...)
}

func (node *ChatNode) println(str string) {
	node.printf("%s\n", str)
}

func (node *ChatNode) printMessage(message *proto.Message) {
	timeStr := time.Unix(message.Timestamp, 0).Format("02.01.2006 15:04")
	node.printf("\n[%s] %s: %s\n", timeStr, message.Name, message.Text)
}

func (node *ChatNode) printError(errorMsg string) {
	node.printf("\nError: %s\n", errorMsg)
}

func discardLine(reader *bufio.Reader) error {
	for {
		_, isPrefix, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if !isPrefix {
			return nil
		}
	}
}

func (node *ChatNode) Run() error {
	return node.RunWithContext(context.Background())
}

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
				err = discardLine(node.reader)
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
