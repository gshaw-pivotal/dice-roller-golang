package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/roll", diceRollHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func diceRollHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have reached the endpoint")
}
