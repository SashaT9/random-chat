package client

import (
	"bufio"
	"net"
	"strconv"
	"time"
)

type Conn struct {
	conn net.Conn
	name string
}

func NewConn(ip string, port int, name string) (*Conn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn: conn,
		name: name,
	}, nil
}

func (c *Conn) WriteMessage(message string) (int, error) {
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return c.conn.Write([]byte(message + "\n"))
}

func (c *Conn) ReadMessage() (string, error) {
	reader := bufio.NewReader(c.conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return message, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
