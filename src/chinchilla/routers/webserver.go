package main

import (
	"chinchilla/loadtest"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var Connection websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ChinchillaPort string

var T *loadtest.LoadTester

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {

	args := os.Args

	if len(args) != 3 {
		fmt.Println("usage is <portno> <chinchilla_portno> for web server")
		os.Exit(1)
	}

	portno := strings.Join([]string{":", args[1]}, "")
	ChinchillaPort = args[2]

	r := mux.NewRouter()

	r.HandleFunc("/", serveIndex).Methods("get")
	r.HandleFunc("/start_test", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		one, _ := strconv.Atoi(r.Form["type_1"][0])
		two, _ := strconv.Atoi(r.Form["type_2"][0])
		three, _ := strconv.Atoi(r.Form["type_3"][0])
		algType := r.Form["alg_type"][0]
		workerCount, _ := strconv.Atoi(r.Form["worker_count"][0])

		p := loadtest.TestParams{one, two, three, algType, workerCount}

		T.Stop = false
		T.StopTest()
		T.LoadTest(w, r, &p)

	}).Methods("post")

	r.HandleFunc("/blast", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		blastType, _ := strconv.Atoi(r.Form["blast_type"][0])
		blastCount, _ := strconv.Atoi(r.Form["blast_count"][0])
		T.Blast(w, r, blastType, blastCount)
	}).Methods("post")

	r.HandleFunc("/stop_test", func(w http.ResponseWriter, r *http.Request) {
		T.Stop = true
		T.StopTest()
	})

	r.HandleFunc("/kill/{wid}", func(w http.ResponseWriter, r *http.Request) {
		wid, _ := strconv.Atoi(mux.Vars(r)["wid"])
		T.KillWorker(w, r, wid)
	}).Methods("post")

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		T.AddWorker()
	}).Methods("post")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			println(err)
			return
		}
		Connection = *conn
		fmt.Println("connected")
		T = loadtest.NewLoadTester(ChinchillaPort, Connection)
	})

	r.PathPrefix("/javascripts/").Handler(http.StripPrefix("/javascripts/", http.FileServer(http.Dir("javascripts/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))

	http.Handle("/", r)
	http.ListenAndServe(portno, nil)

}
