package main

import (
	"fmt"
	"net"

	protocol "github.com/SashaT9/chat-app/pkg"
)

func main() {
	listenner, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer listenner.Close()

	for {
		conn, _ := listenner.Accept()
		fmt.Println("Client connected:", conn.RemoteAddr())
		go func(c net.Conn) {
			defer c.Close()

			buffer := make([]byte, 1024)
			for {
				n, err := c.Read(buffer)
				if err != nil {
					fmt.Println("closed ", err)
					return
				}
				received := string(buffer[:n])
				fmt.Println("Received:", received)
				env, err := protocol.Unmarshal(buffer[:n])
				if err != nil {
					fmt.Println("Error unmarshalling:", err)
					continue
				}
				fmt.Println("Envelope Type:", env.Type)
				responseBytes, err := protocol.RegionCount(42)
				if err != nil {
					fmt.Println("Error marshalling response:", err)
					continue
				}
				responseBytes = append(responseBytes, '\n')
				_, err = c.Write(responseBytes)
				if err != nil {
					fmt.Println("Error writing response:", err)
					return
				}
				fmt.Println("Sent response:", string(responseBytes))
			}
		}(conn)
	}
}
