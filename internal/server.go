package internal

import (
	"bufio"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	Pool []*Client
	Mtx  sync.Mutex
	IP   string
	Port int
}

func Run(ip string, port int) (*Server, error) {
	server := &Server{IP: ip, Port: port}
	listener, err := net.Listen("tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	server.Receive(listener)
	return server, nil
}

func (s *Server) Receive(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
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
			client.Partner.Inbox <- client.Username + ": " + msg + "\n"
		} else {
			client.Inbox <- "No partner connected.\n"
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
	if len(s.Pool) == 0 {
		s.Pool = append(s.Pool, client)
	} else {
		partner := s.Pool[0]
		partner.Partner = client
		client.Partner = partner
		s.Pool = s.Pool[1:]
	}
}

func (s *Server) removeClient(client *Client) {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	if client.Partner != nil {
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
}
