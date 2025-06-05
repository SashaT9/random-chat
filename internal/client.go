package internal

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	Username   string
	Connection net.Conn
	Inbox      chan string
	Partner    *Client
}

func NewClient(username, ip string, port int) (*Client, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	conn.Write([]byte(username + "\n"))

	client := &Client{
		Username:   strings.TrimSpace(username),
		Connection: conn,
		Inbox:      make(chan string, 10),
	}
	go client.Listen()
	return client, nil
}

func (c *Client) Listen() {
	reader := bufio.NewReader(c.Connection)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			close(c.Inbox)
			return
		}
		c.Inbox <- strings.TrimRight(msg, "\n")
	}
}

func (c *Client) Send() error {
	stdinScanner := bufio.NewScanner(os.Stdin)
	if stdinScanner.Scan() {
		line := stdinScanner.Text()
		_, err := c.Connection.Write([]byte(line + "\n"))
		return err
	}
	return stdinScanner.Err()
}

func (c *Client) Read() {
	for msg := range c.Inbox {
		println(msg)
	}
}
