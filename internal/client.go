package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
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
		log.Printf("%s", color.RedString("Error connecting to server: %v\n", err))
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
		msg = strings.TrimRight(msg, "\n")
		c.Inbox <- msg
	}
}

func (c *Client) Send() error {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimRight(line, "\n")
	_, err = c.Connection.Write([]byte(line + "\n"))
	coloredLine := color.GreenString(c.Username+"(you): ") + line
	fmt.Print("\033[A\r\033[K" + coloredLine + "\n")
	return err
}

func (c *Client) Read() {
	for msg := range c.Inbox {
		fmt.Println(msg)
	}
}
