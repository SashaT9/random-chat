package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SashaT9/random-chat/pkg/server"
)

func main() {
	ip := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	fmt.Println(port)
	server.Run(ip, port)
}
