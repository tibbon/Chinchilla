package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"strings"
)

type queue struct {
	conn net.Conn
	qVal uint32
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Println("usage is <portno http> <portno tcp>")
		os.Exit(1)
	}

	portno := strings.Join([]string{":", args[1]}, "")

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "You got root")
	}).Methods("get")

	// Place rest of routes here

	go AcceptWorkers()

	http.Handle("/", r)
	http.ListenAndServe(portno, nil)
}

func AcceptWorkers() {
	portno := strings.Join([]string{":", os.Args[2]}, "")

	ln, err := net.Listen("tcp", portno)
	checkError(err)

	workers := make(map[uint32]queue)
	RespQueue := make(chan mssg.WorkResp)

	for {
		fmt.Println("in loop")
		conn, err := ln.Accept()
		if err != nil {
			continue
		} else {
			fmt.Println("got here")
			go RecvWork(conn, &workers, RespQueue)
		}
	}
}

func RecvWork(conn net.Conn, workers *map[uint32]queue, RespQueue chan mssg.RespQueue) {

	header := new(mssg.Msg)
	resp := new(mssg.WorkResp)
	dec := gob.NewDecoder(conn)
	avgTimes := make(map[uint32]uint32)

	dec.Decode(header)

	if header.Type == 1 && header.Id {
		workers[header.Id] = conn
		fmt.Print("Added slave connection")
	} else {
		conn.Close()
		fmt.Println("improper connect")
		return
	}

	for {
		dec.Decode(resp)
		if resp.Type == 1 {
			conn.Close()
			delete(workers, resp.Id)
		}
	}

}
