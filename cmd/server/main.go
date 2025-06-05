package main

import (
	"fmt"

	"github.com/SashaT9/random-chat/internal"
)

func main() {
	_, err := internal.Run("127.0.0.1", 9000)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
