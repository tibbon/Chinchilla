package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {

	args := os.Args

	if len(args) != 3 {
		fmt.Println("usage is <ip:port> <id>")
		os.Exit(1)
	}
	id, _ := strconv.Atoi(os.Args[2])
	data := mssg.Connect{1, uint32(id), 0}
	wReq := new(mssg.WorkReq)

	if len(args) != 3 {
		fmt.Println("usage is <ip:port> <id>")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", args[1])
	checkError(err)

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	data_struct := new(mssg.WorkRespData)

	enc.Encode(data)
	for {
		err := dec.Decode(wReq)
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("type %u, arg1 %s, host %s\n", wReq.Type, wReq.Arg1, wReq.WId)
		wResp := mssg.WorkResp{1, 1, *data_struct, wReq.WId, 10}
		enc.Encode(wResp)
	}

}
