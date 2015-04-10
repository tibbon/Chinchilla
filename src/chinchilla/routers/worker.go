package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("usage is <ip:port>")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", args[1])
	checkError(err)

	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println("got here")
	fmt.Println(string(result))

}
