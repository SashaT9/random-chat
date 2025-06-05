package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/SashaT9/random-chat/pkg/protocol"
)

type Client struct {
	Username string
	conn     net.Conn
}

func NewClient(username string, ip string, port int) (*Client, error) {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	regMsg := protocol.Message{
		Type:    protocol.Register,
		PayLoad: username,
	}
	fmt.Fprintf(conn, "%s\n", protocol.Marshal(regMsg))
	client := &Client{
		Username: username,
		conn:     conn,
	}
	return client, nil
}

func (c *Client) SendMessage(data string) error {
	m := protocol.Message{
		Type:    protocol.Msg,
		PayLoad: data,
	}
	line := protocol.Marshal(m) + "\n"
	_, err := c.conn.Write([]byte(line))
	return err
}

func (c *Client) ReadLoop() {
	for {
	}
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		raw := scanner.Text()
		msg, err := protocol.Unmarshal(raw)
		if err != nil {
			fmt.Printf("error unmarhalling: %v\n", err)
			continue
		}
		switch msg.Type {
		case protocol.Matched:
			fmt.Printf("You are matched with: %s\n", msg.PayLoad)
		case protocol.Unmatched:
			fmt.Printf("you are unmatched with: %s\n", msg.PayLoad)
		}
	}
}
