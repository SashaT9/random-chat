package main

import (
	"log"
	"os"
	"strconv"

	"github.com/SashaT9/random-chat/internal"
)

func main() {
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	internal.Run(ip, port)
}
