package main

import (
	"chinchilla/mssg"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const poolSize = 15

var Counter int

type Queue struct {
	QVal uint32
	Enc  *gob.Encoder
	Sent bool
	Reqs []mssg.WorkReq
}
type Job struct {
	W   http.ResponseWriter
	Mtx *sync.Mutex
}

type MapJ struct {
	m map[uint32]Job
	l *sync.RWMutex
}

type MapQ struct {
	m map[uint32]Queue
	l *sync.RWMutex
}

type Stack struct {
	s []uint32
	l *sync.RWMutex
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	args := os.Args
	Counter = 0

	if len(args) != 3 {
		fmt.Println("usage is <portno http> <portno tcp>")
		os.Exit(1)
	}
	portno := strings.Join([]string{":", args[1]}, "")

	ReqQueue := make(chan mssg.WorkReq, 10000)
	r := mux.NewRouter()

	jobs := &MapJ{make(map[uint32]Job), new(sync.RWMutex)}
	// nodeQs := make(map[uint32][])
	ids := &Stack{make([]uint32, 10000), new(sync.RWMutex)} // make extensible later

	for i := 0; i < 10000; i++ {
		ids.s[i] = uint32(i)
	}

	r.HandleFunc("/api/{type}/{arg1}", func(w http.ResponseWriter, r *http.Request) {
		var id uint32
		typ, _ := strconv.Atoi(mux.Vars(r)["type"])
		ids.l.Lock()
		id, ids.s = ids.s[0], ids.s[1:] // get a free work id (ultimately this is load distribution)
		ids.l.Unlock()
		AddReqQueue(w, ReqQueue, typ, mux.Vars(r)["arg1"], id, jobs)
		jobs.m[id].Mtx.Lock() // potential map corruption?
	}).Methods("get")

	go AcceptWorkers(ReqQueue, jobs)

	http.Handle("/", r)
	http.ListenAndServe(portno, nil)
}

func AcceptWorkers(ReqQueue chan mssg.WorkReq, jobs *MapJ) {
	portno := strings.Join([]string{":", os.Args[2]}, "")

	ln, err := net.Listen("tcp", portno)
	checkError(err)

	workers := &MapQ{make(map[uint32]Queue), new(sync.RWMutex)}
	RespQueue := make(chan mssg.WorkResp)

	go SendWorkReq(ReqQueue, workers)
	// Makes response pool for compete work requests
	for i := 0; i < poolSize; i++ {
		go SendResp(RespQueue, jobs)
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

func RecvWork(conn net.Conn, workers *MapQ, RespQueue chan mssg.WorkResp) {

	header := new(mssg.Connect)
	resp := new(mssg.WorkResp)
	gob.Register(mssg.WorkReq{})
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	avgTimes := make(map[uint8]int64)
	dec.Decode(header)

	if header.Type == 1 && header.Id != 0 {
		workers.l.Lock()
		workers.m[header.Id] = Queue{header.QVal, enc, false, make([]mssg.WorkReq, 0)}
		workers.l.Unlock()

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
			return
		}
		if resp.Type == 0 {
			conn.Close()
			workers.l.Lock()
			delete(workers.m, resp.Id)
			workers.l.Unlock()
			return
		} else {
			RespQueue <- *resp
			t := time.Since(resp.RTime)
			workers.l.Lock()
			fmt.Printf("Queue length for %d is %d\n", resp.Id, len(workers.m[resp.Id].Reqs))
			tmp := workers.m[resp.Id]
			if len(workers.m[resp.Id].Reqs) != 0 {
				tmp.Reqs = workers.m[resp.Id].Reqs[1:]
				workers.m[resp.Id] = tmp
			}
			workers.l.Unlock()
			avgTimes[resp.Type] = time.Duration.Nanoseconds(t) // Add weighted avg function
		}
	}
}

// Thread to send responses back to hosts
func SendResp(RespQueue chan mssg.WorkResp, jobs *MapJ) {
	for {
		resp := <-RespQueue

		json_resp, _ := json.Marshal(resp)
		// fmt.Println(resp.Data)
		jobs.l.Lock()
		w := jobs.m[resp.WId].W
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)

		// allow cross domain AJAX requests

		_, err := w.Write(json_resp)
		jobs.m[resp.WId].Mtx.Unlock()
		jobs.l.Unlock()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}

	}
}

// Add req struct to a channel
func AddReqQueue(w http.ResponseWriter, ReqQueue chan mssg.WorkReq, typ int, arg1 string, id uint32, jobs *MapJ) {
	jobs.l.Lock()
	jobs.m[id] = Job{W: w, Mtx: &sync.Mutex{}}
	Counter += 1
	fmt.Println(Counter)
	jobs.m[id].Mtx.Lock()
	jobs.l.Unlock()
	ReqQueue <- mssg.WorkReq{Type: uint8(typ), Arg1: arg1, WId: id, STime: time.Now()}
}

func SendWorkReq(ReqQueue chan mssg.WorkReq, workers *MapQ) {
	for {
		req := <-ReqQueue
		node := RoundRobin(workers)
		// fmt.Printf("Sending work request to node %u", node)
		workers.l.Lock()
		tmp := workers.m[node]
		tmp.Reqs = append(workers.m[node].Reqs, req)
		workers.m[node] = tmp
		err := workers.m[node].Enc.Encode(req)
		workers.l.Unlock()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		}
		// fmt.Println("Sent work Request")
		// for k, v := range workers.m {
		// 	workers.l.Lock()
		// 	fmt.Printf("Worker %u", k)
		// 	for i := 0; i < len(v.Reqs); i++ {
		// 		fmt.Printf("%s ", v.Reqs[i].Arg1)
		// 	}
		// 	fmt.Println("")
		// 	workers.l.Unlock()
		// }

	}
}

func RoundRobin(workers *MapQ) uint32 {

	for k, v := range workers.m {
		workers.l.Lock()
		if !workers.m[k].Sent {
			v.Sent = true
			workers.m[k] = v
			workers.l.Unlock()
			return k
		}
		workers.l.Unlock()
	}
	for k, v := range workers.m {
		workers.l.Lock()
		v.Sent = false
		workers.m[k] = v
		workers.l.Unlock()
	}
	for k, v := range workers.m {
		v.Sent = true
		workers.m[k] = v
		return k
	}
	return 0
}
