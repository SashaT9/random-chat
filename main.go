package main

import (
	"fmt"
	"os"

	"github.com/SashaT9/chat-app/internal/client"
)

func main() {
	if len(os.Args) == 1 {
		client.Repl()
		return
	} else {
		fmt.Print("todo")
	}
}
