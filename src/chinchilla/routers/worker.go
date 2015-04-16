package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"fmt"
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
	data := mssg.Connect{1, 1, 0}
	wReq := new(mssg.WorkReq)

	if len(args) != 2 {
		fmt.Println("usage is <ip:port>")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", args[1])
	checkError(err)

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	enc.Encode(data)
	for {
		err = dec.Decode(wReq)
		if err != nil {
			fmt.Println("wrong send 1?")
		}
		fmt.Printf("type %u, arg1 %s, host %s\n", wReq.Type, wReq.Arg1, wReq.WId)
		wResp := mssg.WorkResp{1, 1, []byte("You win my the data!"), wReq.WId, 10}
		err = enc.Encode(wResp)
		if err != nil {
			fmt.Println("wrong send 2?")
		}
	}

}
