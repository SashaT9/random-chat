package internal

import (
	"bufio"
	"log"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Server struct {
	Pool []*Client
	Mtx  sync.Mutex
	IP   string
	Port int
}

func Run(ip string, port int) error {
	log.Printf("Starting server on %s:%d\n", ip, port)
	server := &Server{IP: ip, Port: port}
	listener, err := net.Listen("tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		return err
	}
	server.Receive(listener)
	return nil
}

func (s *Server) Receive(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		r := bufio.NewReader(conn)
		username, _ := r.ReadString('\n')
		username = strings.TrimSpace(username)
		client := &Client{
			Username:   username,
			Connection: conn,
			Inbox:      make(chan string, 10),
			Partner:    nil,
		}
		s.addClient(client)
		go s.HandleRead(client, r)
		go s.HandleWrite(client)
	}
}

func (s *Server) HandleRead(client *Client, r *bufio.Reader) {
	defer func() {
		s.removeClient(client)
		close(client.Inbox)
		client.Connection.Close()
	}()
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.TrimSpace(msg)
		if client.Partner != nil {
			client.Partner.Inbox <- color.BlueString("%s: ", client.Username) + msg + "\n"
		} else {
			go func() {
				client.Inbox <- color.RedString("Waiting for a partner...") + "\n"
				for client.Partner == nil {
					time.Sleep(10 * time.Millisecond)
				}
				client.Partner.Inbox <- color.BlueString("%s: ", client.Username) + msg + "\n"
			}()
		}
	}
}

func (s *Server) HandleWrite(client *Client) {
	writer := bufio.NewWriter(client.Connection)
	for msg := range client.Inbox {
		writer.WriteString(msg)
		writer.Flush()
	}
}

func (s *Server) addClient(client *Client) {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	s.Pool = append(s.Pool, client)
	s.updatePool()
	log.Printf("Client %s connected\n", client.Username)
}

func (s *Server) removeClient(client *Client) {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	if client.Partner != nil {
		client.Partner.Inbox <- color.RedString("Your partner %s disconnected.", client.Username) + "\n"
		client.Partner.Partner = nil
		s.Pool = append(s.Pool, client.Partner)
		client.Partner = nil
	}
	for i, c := range s.Pool {
		if c == client {
			s.Pool = slices.Delete(s.Pool, i, i+1)
			break
		}
	}
	s.updatePool()
	log.Printf("Client %s disconnected\n", client.Username)
}

func (s *Server) updatePool() {
	for len(s.Pool) > 1 {
		client1 := s.Pool[0]
		client2 := s.Pool[1]
		client1.Inbox <- color.CyanString("Connected to: ") + client2.Username + "\n"
		client2.Inbox <- color.CyanString("Connected to: ") + client1.Username + "\n"
		client1.Partner = client2
		client2.Partner = client1
		s.Pool = s.Pool[2:]
	}
}
