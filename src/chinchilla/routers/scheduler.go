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
	fmt.Println(portno)
	ln, err := net.Listen("tcp", portno)
	checkError(err)
	for {
		fmt.Println("in loop")
		conn, err := ln.Accept()
		if err != nil {
			continue
		} else {
			fmt.Println("got here")
			go RecvWork(conn)
		}

	}
}

func RecvWork(conn net.Conn) {
	data := new(mssg.Connect)
	dec := gob.NewDecoder(conn)
	dec.Decode(data)
	fmt.Printf("op %d, id % d", data.Op, data.Id)
}
