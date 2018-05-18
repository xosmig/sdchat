package sdchat

import (
	"bufio"
	"fmt"
	"github.com/xosmig/sdchat2/proto"
	"time"
	"github.com/xosmig/sdchat2/apiclient"
	"os"
	"log"
	"math/rand"
	"io"
)

type ChatNode struct {
	reader    *bufio.Reader
	name      string
	apiClient apiclient.ApiClient
}

func NewChatNode(name string, apiClient apiclient.ApiClient) ChatNode {
	return ChatNode{bufio.NewReader(os.Stdin), name, apiClient}
}

func (node *ChatNode) sendMessage(text string) {
	message := &sdchat2.Message{
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
	fmt.Printf(format, a...)
}

func (node *ChatNode) println(str string) {
	node.printf("%s\n", str)
}

func (node *ChatNode) printMessage(message *sdchat2.Message) {
	timeStr := time.Unix(message.Timestamp, 0).Format("02.01.2006 15:04")
	node.printf("\n[%s] %s: %s\n", timeStr, message.Name, message.Text)
}

func (node *ChatNode) printError(errorMsg string) {
	node.printf("\nError: %s\n", errorMsg)
}

func (node *ChatNode) sayBye() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	bye := []string{"Tschüss!", "Bye!", "Пока!", "バイ!"}
	node.println(bye[rnd.Intn(len(bye))])
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
	log.Println("Connecting...")
	err := node.apiClient.Start()
	if err != nil {
		node.printError(fmt.Sprintf("Connection error: %v", err))
		return err
	}
	defer func() { node.apiClient.Stop() }()
	log.Println("Connected.")
	log.Println("Type \"m[enter]\" to write a message or \"q[enter]\" to exit")

	finish := make(chan bool, 1)
	defer func() { close(finish) }()

	lines := make(chan string)
	go func() {
		for {
			bytes, isPrefix, err := node.reader.ReadLine()
			if isPrefix {
				node.printError("Line is too long")
				err = discardLine(node.reader)
			}
			if err == io.EOF {
				finish <- true
				break
			}
			if err != nil {
				node.printError(fmt.Sprintf("Oops: error reading your input: %v", err))
				continue
			}
			lines <- string(bytes)
		}
	}()

	messages := make(chan *sdchat2.Message)
	defer func() { close(messages) }()
	go func() {
		for {
			message, err := node.apiClient.ReceiveMessage()
			if err != nil {
				node.println("Connection is closed")
				finish <- true
				break
			}
			messages <- message
		}
	}()

selectLoop:
	for {
		select {
		case <-finish:
			break selectLoop
		case line := <-lines:
			switch line {
			case "m":
				node.printf("Enter message: ")
				text := <-lines  // blocking
				node.sendMessage(text)
			case "q":
				finish <- true
			default:
				node.printError(fmt.Sprintf("Unknown command: '%s'", line))
			}
		case message := <-messages:
			node.printMessage(message)
		}
	}

	node.sayBye()
	return nil
}
