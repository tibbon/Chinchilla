package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	data := &mssg.Connect{0, 234}

	if len(args) != 2 {
		fmt.Println("usage is <ip:port>")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", args[1])
	checkError(err)
	enc := gob.NewEncoder(conn)

	enc.Encode(data)

	time.Sleep(100000 * time.Millisecond)

}
