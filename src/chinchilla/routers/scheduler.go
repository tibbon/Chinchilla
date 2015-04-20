package main

import (
	"chinchilla/mssg"
	"chinchilla/send"
	"chinchilla/types"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const poolSize = 15

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

	ReqQueue := make(chan mssg.WorkReq, 10000)
	r := mux.NewRouter()

	jobs := &types.MapJ{make(map[uint32]types.Job), new(sync.RWMutex)}
	ids := &types.Stack{make([]uint32, 10000), new(sync.RWMutex)} // make extensible later

	for i := 0; i < 10000; i++ {
		ids.S[i] = uint32(i)
	}

	r.HandleFunc("/api/{type}/{arg1}", func(w http.ResponseWriter, r *http.Request) {
		var id uint32
		typ, _ := strconv.Atoi(mux.Vars(r)["type"])
		ids.L.Lock()
		id, ids.S = ids.S[0], ids.S[1:] // get a free work id (ultimately this is load distribution)
		ids.L.Unlock()
		send.Scheduler(w, ReqQueue, typ, mux.Vars(r)["arg1"], id, jobs)
		fmt.Printf("got here with id %d\n", id)
		<-jobs.M[id].Sem

	}).Methods("get")

	go AcceptWorkers(ReqQueue, jobs)

	http.Handle("/", r)
	http.ListenAndServe(portno, nil)
}

func AcceptWorkers(ReqQueue chan mssg.WorkReq, jobs *types.MapJ) {
	portno := strings.Join([]string{":", os.Args[2]}, "")

	ln, err := net.Listen("tcp", portno)
	checkError(err)

	workers := &types.MapQ{make(map[uint32]types.Queue), new(sync.RWMutex)}
	RespQueue := make(chan mssg.WorkResp)

	go send.Node(ReqQueue, workers)
	// Makes response pool for compete work requests
	for i := 0; i < poolSize; i++ {
		go send.Client(RespQueue, jobs)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		} else {
			go RecvWork(conn, workers, RespQueue)
		}
	}
}

func RecvWork(conn net.Conn, workers *types.MapQ, RespQueue chan mssg.WorkResp) {

	header := new(mssg.Connect)
	resp := new(mssg.WorkResp)
	gob.Register(mssg.WorkReq{})
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	dec.Decode(header)

	if header.Type == 1 && header.Id != 0 {
		workers.L.Lock()
		workers.M[header.Id] = types.Queue{header.QVal, enc, false, make([]mssg.WorkReq, 0), make(map[uint8]float64)}
		workers.L.Unlock()

	} else {
		conn.Close()
		fmt.Println("improper connect")
		return
	}
	// Loop until server send 1 (D/C) or process infinite responses and update time objects and add to queue
	for {
		err := dec.Decode(resp)
		if err != nil {
			conn.Close()
			workers.L.Lock()
			delete(workers.M, resp.Id)
			workers.L.Unlock()
			return
		}
		if resp.Type == 0 {
			conn.Close()
			workers.L.Lock()
			delete(workers.M, resp.Id)
			workers.L.Unlock()
			return
		} else {
			RespQueue <- *resp
			HandleResp(resp, workers, header.Id)
		}
	}
}

func HandleResp(resp *mssg.WorkResp, workers *types.MapQ, id uint32) {
	t := resp.RTime
	fmt.Printf("Queue length for %d is %d\n", resp.Id, len(workers.M[resp.Id].Reqs))

	workers.L.Lock()
	tmp := workers.M[resp.Id]

	if len(workers.M[resp.Id].Reqs) != 0 {
		tmp.Reqs = workers.M[resp.Id].Reqs[1:]
		workers.M[resp.Id] = tmp
	}

	workers.M[id].AvgTimes[resp.Type] = t
	tmp = workers.M[id]

	tmp.QVal -= t - .1
	if tmp.QVal < 0.0 {
		tmp.QVal = 0
	}
	workers.M[id] = tmp

	workers.L.Unlock()
	fmt.Printf("time is %f", t)
}
