package nativemessage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
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

func fromJson(j string) (interface{}, error) {
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

func toJson(m interface{}) ([]byte, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return j, nil
}

type NativeMessageReader struct {
	ready     bool
	byteOrder binary.ByteOrder
}

func NewNativeMessageReader() NativeMessageReader {
	byteOrder := getByteOrder()

	n := NativeMessageReader{ready: true, byteOrder: byteOrder}

	return n
}

func getByteOrder() binary.ByteOrder {
	var one int16 = 1
	var byteOrder binary.ByteOrder
	b := (*byte)(unsafe.Pointer(&one))
	if *b == 0 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
	return byteOrder
}

func (n *NativeMessageReader) ReadMessage() (interface{}, error) {
	// Read message size
	var messageLength int32
	err := binary.Read(os.Stdin, n.byteOrder, &messageLength)
	if err != nil {
		return nil, err
	}
	log.Printf("received message size: %d", messageLength)
	// Read message bytes
	message := make([]byte, messageLength)
	_, err = io.ReadAtLeast(os.Stdin, message, int(messageLength))
	// r = bufio.NewReaderSize(os.Stdin, int(messageLength))
	// _, err = r.Read(message)
	if err != nil {
		return nil, err
	}
	// Parse and return message
	log.Printf("received message size: %s", message)
	return fromJson(string(message))
}

func EncodeNativeMessage(m interface{}) ([]byte, error) {
	buff := new(bytes.Buffer)
	nm, err := toJson(m)
	if err != nil {
		return nil, err
	}
	byteOrder := getByteOrder()
	binary.Write(buff, byteOrder, uint32(len(nm)))
	binary.Write(buff, byteOrder, nm)
	return buff.Bytes(), nil
}
