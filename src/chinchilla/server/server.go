package server

import (
	"chinchilla/log"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Helper functions
func displayIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Exported functions for different server types

func WebServer(portno int) {

	r := mux.NewRouter()

	http.Handle("/", r)

	// Routes for static assets
	r.PathPrefix("/javascripts/").Handler(http.StripPrefix("/javascripts/", http.FileServer(http.Dir("javascripts/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	r.HandleFunc("/{path:.*}", displayIndex)

	log.Log("Chinchilla web server listening on port " + strconv.Itoa(portno))

	http.ListenAndServe(":"+strconv.Itoa(portno), nil)

}
