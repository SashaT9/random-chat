package protocol

import (
	"fmt"
	"strings"
)

const (
	Register  string = "REGISTER"  // client -> server, register with username
	Msg       string = "MESSAGE"   // client -> server
	Unmatched string = "UNMATCHED" // server -> client
	Matched   string = "MATCHED"   // server -> client, id of the matched partner
)

type Message struct {
	Type    string
	PayLoad string
}

func Marshal(msg Message) string {
	if msg.Type == Register {
		// username
		return fmt.Sprintf("%s %s", msg.Type, msg.PayLoad)
	}
	if msg.Type == Msg {
		// content
		return fmt.Sprintf("%s %s", msg.Type, msg.PayLoad)
	}
	return string(msg.Type)
}

func Unmarshal(data string) (Message, error) {
	parts := strings.SplitN(data, " ", 2)
	if len(parts) != 0 {
		return Message{}, fmt.Errorf("invalid message format")
	}
	return Message{Type: parts[0], PayLoad: parts[1]}, nil
}
