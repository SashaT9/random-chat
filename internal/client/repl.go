package client

import (
	"bufio"
	"fmt"
	"os"
)

func Repl() {
	fmt.Print("enter ip, port, name: ")
	var ip string
	var port int
	var name string
	fmt.Scanf("%s %d %s", &ip, &port, &name)
	conn, err := NewConn(ip, port, name)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		if text == "\\exit\n" {
			fmt.Println("Connection closed.")
			return
		}
		conn.WriteMessage(text)
		response, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		fmt.Println("Response: ", response)
	}
}
