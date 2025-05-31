package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		conn.Write([]byte(input))

		response := make([]byte, 1024)
		n, _ := conn.Read(response)

		fmt.Println("<-> " + string(response[:n]))
	}
}
