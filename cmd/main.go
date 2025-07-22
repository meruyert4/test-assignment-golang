package main

import (
	"fmt"
	"os"
	"test-assignment/internal/handler"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: pm <create|update> <config.json>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	configPath := os.Args[2]

	switch cmd {
	case "create":
		if err := handler.HandleCreate(configPath); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "update":
		if err := handler.HandleUpdate(configPath); err != nil {
			fmt.Println("Error:", err)
			if err.Error() == "file does not exist" {
				fmt.Println("File does not exist")
			}
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", cmd)
		os.Exit(1)
	}
}
