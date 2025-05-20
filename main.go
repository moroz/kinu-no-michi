package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const LISTEN_ON = ":3000"

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<h1>Hello world!</h1>`)
}

func main() {
	r := chi.NewRouter()

	r.Get("/", handleIndex)

	log.Printf("Listening on %s", LISTEN_ON)

	log.Fatal(http.ListenAndServe(LISTEN_ON, r))
}
