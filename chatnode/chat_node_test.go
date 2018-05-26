package sdchat

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bouk/monkey"
	"github.com/golang/mock/gomock"
	"github.com/xosmig/sdchat/apiclient/mock_apiclient"
	"github.com/xosmig/sdchat/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestChatNode(t *testing.T) {
	// stub time.Now with a deterministic value
	patch := monkey.Patch(time.Now, func() time.Time {
		return time.Date(2018, time.May, 20, 23, 59, 0, 0, time.Local)
	})
	defer patch.Unpatch()

	type message struct {
		text string
		sync bool
	}
	EXIT := ":EXIT:"
	WAIT := ":WAIT:"
	testCases := []struct {
		name    string
		msgs    []message
		replies []message
	}{
		{"empty", []message{{EXIT, false}}, []message{{WAIT, false}}},
		{"simpleGreeting",
			[]message{{"hi!", true}, {EXIT, false}},
			[]message{{"hallo", true}, {WAIT, false}}},
		{"twoMessagesInARow",
			[]message{{"Hi! How are you?", true}, {"Magnificent!", false}, {"Bye!", true}, {EXIT, false}},
			[]message{{"Hi, great!", true}, {"c u", true}, {WAIT, false}}},
		{"twoRepliesInARow",
			[]message{{"Hi! How are you?", true}, {WAIT, false}},
			[]message{{"Hi", true}, {"Don't wanna talk", false}, {EXIT, false}}},
	}

	log.SetOutput(ioutil.Discard)

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			wantReply := make(chan bool)

			mockApiClient := mock_apiclient.NewMockApiClient(ctrl)
			mockApiClient.EXPECT().Start().Times(1).Return(nil)
			for _, replyIter := range tc.replies {
				reply := replyIter
				mockApiClient.EXPECT().ReceiveMessage().Times(1).DoAndReturn(func() (*proto.Message, error) {
					if reply.sync {
						<-wantReply
						time.Sleep(100 * time.Millisecond)
					}
					if reply.text == EXIT {
						return nil, fmt.Errorf("canceled connection")
					}
					if reply.text == WAIT {
						time.Sleep(1 * time.Second)
					}
					message := &proto.Message{
						Name:      "Mock",
						Text:      reply.text,
						Timestamp: time.Now().Unix(),
					}
					return message, nil
				})
			}
			mockApiClient.EXPECT().SendMessage(gomock.Any()).Times(len(tc.msgs) - 1).Return(nil)
			mockApiClient.EXPECT().Stop().Times(1)

			chatNode := NewChatNode("test", mockApiClient)

			output := new(bytes.Buffer)
			chatNode.stdout = output

			inputReader, inputWriter := io.Pipe()
			chatNode.reader = bufio.NewReader(inputReader)

			done := make(chan bool)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				chatNode.RunWithContext(ctx)
				done <- true
			}()

			for _, msg := range tc.msgs {
				if msg.text == WAIT {
					time.Sleep(1 * time.Second)
				} else if msg.text == EXIT {
					fmt.Fprintln(inputWriter, "q")
				} else {
					fmt.Fprintln(inputWriter, "m")
					fmt.Fprintln(inputWriter, msg.text)
				}
				if msg.sync {
					wantReply <- true
					time.Sleep(200 * time.Millisecond)
				}
			}

			select {
			case <-done:
				// OK, continue
			case <-time.After(1 * time.Second):
				tt.Fatalf("Timeout")
			}

			answerFile := fmt.Sprintf("./tests_output/TestChatNode/%s", tc.name)
			if os.Getenv("SDCHAT_UPDATE_TEST_ANSWER") != "" {
				err := ioutil.WriteFile(answerFile, output.Bytes(), 0666)
				if err != nil {
					t.Fatalf("Error updating test answer")
				}
			} else {
				expected, err := ioutil.ReadFile(answerFile)
				if err != nil {
					tt.Fatalf("Cannot open answer file")
				}
				if string(expected) != string(output.Bytes()) {
					tt.Errorf("Invalid output:\n%v", string(output.Bytes()))
				}
			}
		})
	}
}
