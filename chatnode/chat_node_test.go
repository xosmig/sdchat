package sdchat

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bouk/monkey"
	"github.com/golang/mock/gomock"
	"github.com/xosmig/sdchat/apiclient/mock_apiclient"
	"github.com/xosmig/sdchat/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// special text for mock messages
const EXIT = ":EXIT:"

const SLEEP = ":SLEEP:"
const SYNC = ":SYNC:"
const WAIT = ":WAIT:"

func doSynchronization(msgText *string, syncChan chan struct{}) {
	for {
		if strings.HasPrefix(*msgText, SLEEP) {
			time.Sleep(1 * time.Second)
			*msgText = (*msgText)[len(SLEEP):]
			continue
		}

		if strings.HasPrefix(*msgText, WAIT) {
			<-syncChan
			// this sleep is necessary
			time.Sleep(300 * time.Millisecond)
			*msgText = (*msgText)[len(WAIT):]
			continue
		}

		if strings.HasSuffix(*msgText, SYNC) {
			syncChan <- struct{}{}
			*msgText = (*msgText)[:len(*msgText)-len(SYNC)]
			continue
		}

		break
	}
}

func mockExpectReply(mock *mock_apiclient.MockApiClient, replyMsg string, syncChan chan struct{}) {
	mock.EXPECT().ReceiveMessage().Times(1).DoAndReturn(func() (*proto.Message, error) {
		doSynchronization(&replyMsg, syncChan)

		if replyMsg == EXIT {
			return nil, fmt.Errorf("canceled connection")
		}

		return &proto.Message{
			Name:      "Mock",
			Text:      replyMsg,
			Timestamp: time.Now().Unix(),
		}, nil
	})
}

func TestChatNode(t *testing.T) {
	// stub time.Now with a deterministic value
	patch := monkey.Patch(time.Now, func() time.Time {
		return time.Date(2018, time.May, 20, 23, 59, 0, 0, time.Local)
	})
	defer patch.Unpatch()

	testCases := []struct {
		name    string
		msgs    []string
		replies []string
	}{
		{"empty",
			[]string{EXIT},
			[]string{WAIT + EXIT},
		},
		{"simpleGreeting",
			[]string{"hi!" + SYNC, WAIT + EXIT + SYNC},
			[]string{WAIT + "hallo" + SYNC, WAIT + EXIT},
		},
		{"twoMessagesInARow",
			[]string{"Hi! How are you?" + SYNC, WAIT + "Magnificent!", "Bye!" + SYNC, WAIT + EXIT + SYNC},
			[]string{WAIT + "Hi, great!" + SYNC, WAIT + "c u" + SYNC, WAIT + EXIT},
		},
		{"twoRepliesInARow",
			[]string{"Hi! How are you?" + SYNC, WAIT + EXIT},
			[]string{WAIT + "Hi", "Don't wanna talk", SLEEP + EXIT + SYNC},
		},
	}

	log.SetOutput(ioutil.Discard)

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			syncChan := make(chan struct{})

			// set the mock expectations
			mockApiClient := mock_apiclient.NewMockApiClient(ctrl)
			mockApiClient.EXPECT().Start().Times(1).Return(nil)
			for _, reply := range tc.replies {
				mockExpectReply(mockApiClient, reply, syncChan)
			}
			mockApiClient.EXPECT().SendMessage(gomock.Any()).Times(len(tc.msgs) - 1).Return(nil)
			mockApiClient.EXPECT().Stop().Times(1)

			// create a chat node
			chatNode := NewChatNode("test", mockApiClient)

			// innit IO utility
			output := new(bytes.Buffer)
			chatNode.stdout = output
			inputReader, inputWriter := io.Pipe()
			chatNode.reader = bufio.NewReader(inputReader)

			// start the chat node
			done := make(chan bool)
			go func() {
				chatNode.Run()
				done <- true
			}()

			// emulate client's part
			for _, msg := range tc.msgs {
				doSynchronization(&msg, syncChan)

				if msg == EXIT {
					fmt.Fprintln(inputWriter, "q")
					break
				}

				fmt.Fprintln(inputWriter, "m")
				fmt.Fprintln(inputWriter, msg)
			}

			// wait for chat node completion with 1-second timeout
			select {
			case <-done:
				// OK, continue
			case <-time.After(1 * time.Second):
				tt.Fatalf("Timeout")
			}

			// compare the result with the golden image
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
