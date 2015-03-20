package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Helper functions
func displayIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func masterRequest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Exported functions for different server types

/* Web server routing. Used to display statistics generated
 * by other server types */
func WebServer(port int) {

	r := mux.NewRouter()

	http.Handle("/", r)

	// Routes for static assets
	r.PathPrefix("/javascripts/").Handler(http.StripPrefix("/javascripts/", http.FileServer(http.Dir("javascripts/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	r.HandleFunc("/{path:.*}", displayIndex)

	// log.Log("Chinchilla web server listening on port " + strconv.Itoa(port))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}

func MasterServer(port int) {

	r := mux.NewRouter()

	http.Handle("/", r)

	r.HandleFunc("/", displayIndex).methods("post")

	// log.Log("Chinchilla master server listening on port " + strconv.Itoa(port))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}

func SlaveServer() {

}
