package cli

import (
	"bufio"
	"consistent-hashing/hash"
	"fmt"
	"os"
	"strings"
)

func Cli(hashRing *hash.HashRing) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		var input string
		fmt.Printf("> ")
		if scanner.Scan() {
			input = scanner.Text()
		}
		inputSlice := strings.Split(input, " ")
		command := inputSlice[0]

		switch command {
		case "q", "quit":
			os.Exit(0)
		case "add":
			hashRing.AddServer(inputSlice[1])
		case "remove":
			hashRing.RemoveServer(inputSlice[1])
		case "ls":
			hashRing.ListServerInfo()
		default:
			fmt.Println("Undefined command!")
		}

	}
}
