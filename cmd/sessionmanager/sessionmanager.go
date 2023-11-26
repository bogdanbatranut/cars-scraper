package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("starting SessionMANAGER service...")

	r := mux.NewRouter().StrictSlash(true)
	log.Println("here")

	r.HandleFunc("/echo", echo()).Methods("Get")
	log.Println("handle")
	err := http.ListenAndServe(":3336", r)
	if err != nil {
		panic(err)
	}

	log.Println("Listening on port 3336")
}

func echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Session MANAGER"))
	}
}
