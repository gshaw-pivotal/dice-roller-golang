package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type DieRollRequest struct {
	Dies []DieRoll `json:"dies"`
}

type DieRoll struct {
	DieType   int `json:"dieType"`
	RollCount int `json:"rollCount"`
}

type DieRollResponse struct {
	DieType     int
	RollCount   int
	RollResults []int
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/roll", diceRollHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func diceRollHandler(w http.ResponseWriter, r *http.Request) {
	var rollList DieRollRequest
	err := json.NewDecoder(r.Body).Decode(&rollList)
	if err != nil {
		fmt.Fprintf(w, "Bad request: %s", err)
		return
	}

	if len(rollList.Dies) > 0 {
		var resultingRoll = make([]DieRollResponse, len(rollList.Dies))
		for index, dieRoll := range rollList.Dies {
			resultingRoll[index] = DieRollResponse{
				DieType:     dieRoll.DieType,
				RollCount:   dieRoll.RollCount,
				RollResults: make([]int, dieRoll.RollCount),
			}

			for rollNum := 0; rollNum < dieRoll.RollCount; rollNum++ {
				resultingRoll[index].RollResults[rollNum] = rand.Intn(dieRoll.DieType) + 1
			}
		}
		json.NewEncoder(w).Encode(resultingRoll)
		return
	}

	fmt.Fprintf(w, "No die rolls specs found in request")
}
