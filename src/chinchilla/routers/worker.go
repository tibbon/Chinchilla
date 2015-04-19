package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
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

	enc.Encode(data)
	rand.Seed(time.Now().Unix())
	for {
		err := dec.Decode(wReq)
		if err != nil {
			conn.Close()
			return
		}
		handleRequest(wReq, enc, id)
	}

	fmt.Println("here")

}

func handleRequest(wReq *mssg.WorkReq, enc *gob.Encoder, wId uint32) {

	data_struct := new(mssg.WorkRespData)
	work_time := 0.0

	switch wReq.Type {
	case 1:
		work_time = (rand.Float64() * 0.5) + 0.05
	case 2:
		work_time = (rand.Float64() * 0.75) + 0.5
	case 3:
		work_time = (rand.Float64() * 1) + 0.75
	}

	fmt.Println(work_time)

	time.Sleep(time.Duration(work_time*1000) * time.Millisecond)

	fmt.Printf("type %u, arg1 %s, host %s\n", wReq.Type, wReq.Arg1, wReq.WId)
	wResp := mssg.WorkResp{1, wId, *data_struct, wReq.WId, work_time}
	enc.Encode(wResp)
}
