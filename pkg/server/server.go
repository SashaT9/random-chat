package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/SashaT9/random-chat/pkg/protocol"
)

type Client struct {
	Username string
	Conn     net.Conn
	Inbox    chan string
}

var (
	pool []*Client
	mtx  sync.Mutex
)

func Run(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed accepting: %v\n", err)
			continue
		}
		fmt.Printf("New connection from %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		conn.Close()
		return
	}
	line := scanner.Text()
	msg, err := protocol.Unmarshal(line)
	if err != nil || msg.Type != protocol.Register {
		conn.Close()
		return
	}
	c := &Client{
		Username: msg.PayLoad,
		Conn:     conn,
		Inbox:    make(chan string, 10),
	}
	go clientWriter(c)

	mtx.Lock()
	if len(pool) == 0 {
		fmt.Printf("New client connected: %s\n", c.Username)
		pool = append(pool, c)
		mtx.Unlock()
	} else {
		partner := pool[0]
		pool = pool[1:]
		fmt.Printf("New client connected: %s, matched with: %s\n", c.Username, partner.Username)
		mtx.Unlock()
	}
}

func clientWriter(c *Client) {
	for msg := range c.Inbox {
		fmt.Fprintln(c.Conn, msg)
	}
}
