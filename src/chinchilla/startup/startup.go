package main

import (
	"chinchilla/server"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	args := os.Args

	if len(args) != 2 {

		fmt.Printf("ERROR: Usage: %s <server_type> <type_args> | server_type are webserver, master, slave\n", args[0])
		return
	}

	server_type := args[1]

	if strings.EqualFold(server_type, "webserver") {

		if len(args) != 3 {

			fmt.Printf("ERROR: Usage: %s %s <portno>\n", args[0], args[1])

		} else {

			portno, err := strconv.Atoi(args[2])
			fmt.Print(err)
			server.WebServer(portno)

		}

	} else if server_type == "master" {
		// server.LoadBalancer()
	}

}
