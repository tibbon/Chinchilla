package main

import (
	"chinchilla/server"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func startServer(args []string) {

	if len(args) != 3 {
		fmt.Printf("ERROR: Usage: %s %s <portno>\n", args[0], args[1])
	} else {

		server_type := args[1]
		portno, err := strconv.Atoi(args[2])
		if err == nil {
			switch server_type {
			case "webserver":
				server.WebServer(portno)
			case "master":
				server.MasterServer(portno)
			case "slave":
				// server.SlaveServer(portno)
			}

		} else {
			fmt.Printf("ERROR: Usage: %s %s <portno>\n", args[0], args[1])
		}

	}
}

func main() {

	args := os.Args

	if len(args) < 2 {

		fmt.Printf("ERROR: Usage: %s <server_type> <type_args> | server_type are webserver, master, slave\n", args[0])
		return
	} else {
		server_type := args[1]
		for _, a := range [3]string{"webserver", "master", "slave"} {
			if strings.EqualFold(server_type, a) {
				startServer(args)
			}
		}
		fmt.Printf("ERROR: Usage: %s <server_type> <type_args> | server_type are webserver, master, slave\n", args[0])
	}

}
