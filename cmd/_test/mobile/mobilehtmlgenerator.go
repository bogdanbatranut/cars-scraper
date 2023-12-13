package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("starting htmlgeneratory service...")

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", getHtml()).Methods("GET")

	appPort := "8989"
	err := http.ListenAndServe(fmt.Sprintf(":%s", appPort), r)
	if err != nil {
		panic(err)
	}
	log.Printf("HTTP listening on port %s\n", appPort)
}

func getHtml() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("cmd/_test/mobile/page.html")
		if err != nil {
			panic(err)
		}
		w.Write(data)

	}
}
