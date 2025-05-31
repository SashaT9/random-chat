package main

import (
	"fmt"
	"net"
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
				response := "echo: " + string(buffer[:n])
				fmt.Println("Received:", response)
				c.Write([]byte(response))
			}
		}(conn)
	}
}
