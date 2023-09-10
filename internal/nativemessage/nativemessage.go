package nativemessage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type message struct {
	Action string `json:"action"`
}

type messageLanguage struct {
	QtGadgetsLanguage int `json:"qtGadgetsLanguage"`
}

type messageFunctionCall struct {
	FunctionName      string `json:"functionName"`
	FunctionArguments string `json:"FunctionArguments"`
}

type MessageStart struct {
	message
	Params messageLanguage `json:"params"`
}

type MessageForward struct {
	message
	Params messageFunctionCall `json:"params"`
}

type QtToZDMessage struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}

func FromJson(j string) (interface{}, error) {
	var m message
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		return nil, err
	}
	switch m.Action {
	case "START":
		var s MessageStart
		if err := json.Unmarshal([]byte(j), &s); err != nil {
			return nil, err
		}
		return s, nil

	case "FORWARD":
		var s MessageForward
		if err := json.Unmarshal([]byte(j), &s); err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, fmt.Errorf("unknown action")
	}
}

func (m message) ToJson() (string, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

type NativeMessageReader struct {
	ready  bool
	reader bufio.Reader
}

func NewNativeMessageReader() NativeMessageReader {
	reader := bufio.NewReader(os.Stdin)
	n := NativeMessageReader{reader: *reader, ready: true}

	return n
}

func (n *NativeMessageReader) ReadMessage() (interface{}, error) {
	s, err := n.reader.ReadString('\n')

	if err != nil {
		log.Println("Error in read string", err)
	}
	// TODO native message marshal
	m, err := FromJson(s)
	if err != nil {
		return nil, err
	}
	return m, nil
}
