package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const poolSize = 15

type Queue struct {
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

	ReqQueue := make(chan mssg.WorkReq)
	r := mux.NewRouter()

	r.HandleFunc("/api/{type}/{arg1}", func(w http.ResponseWriter, r *http.Request) {
		typ, _ := strconv.Atoi(mux.Vars(r)["type"])
		AddReqQueue(w, r, ReqQueue, typ, mux.Vars(r)["arg1"])
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

	workers := make(map[uint32]Queue)
	RespQueue := make(chan mssg.WorkResp)

	// Makes response pool for compete work requests
	for i := 0; i < poolSize; i++ {
		go SendResp(RespQueue)
	}

	for {
		fmt.Println("in loop")
		conn, err := ln.Accept()
		if err != nil {
			continue
		} else {
			fmt.Println("got here")
			go RecvWork(conn, workers, RespQueue)
		}
	}
}

func RecvWork(conn net.Conn, workers map[uint32]Queue, RespQueue chan mssg.WorkResp) {

	header := new(mssg.Connect)
	resp := new(mssg.WorkResp)
	dec := gob.NewDecoder(conn)
	avgTimes := make(map[uint8]uint32)

	dec.Decode(header)

	if header.Type == 1 && header.Id != 0 {
		workers[header.Id] = Queue{conn, header.QVal}
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
			return
		} else {
			RespQueue <- *resp               // May be pointer issue, need to test hard
			avgTimes[resp.Type] = resp.RTime // Add weighted avg function
		}
	}
}

func SendResp(RespQueue chan mssg.WorkResp) {
	for {
		resp := <-RespQueue
		conn, err := net.Dial("tcp", resp.Host)
		if err != nil {
			continue
		}
		conn.Write(resp.Data) // writes raw bytes
	}
}

func AddReqQueue(w http.ResponseWriter, r *http.Request, ReqQueue chan mssg.WorkReq, typ int, arg1 string) {
	// host := r.RemoteAddr // gets  hostname sender
}
