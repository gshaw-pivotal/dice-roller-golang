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

type UsageStatsResponse struct {
	TotalRequests int
	GoodRequests  int
	BadRequests   int
	DieTypeUsage  map[int]int
}

var totalRequestCount int
var goodRequestCount int
var badRequestCount int

var dieTypeMap map[int]int

func main() {
	dieTypeMap = make(map[int]int)

	router := mux.NewRouter()

	router.HandleFunc("/roll", diceRollHandler).Methods("POST")
	router.HandleFunc("/stats", getUsageStatsHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getUsageStatsHandler(w http.ResponseWriter, r *http.Request) {
	var response UsageStatsResponse

	response.TotalRequests = totalRequestCount
	response.GoodRequests = goodRequestCount
	response.BadRequests = badRequestCount
	response.DieTypeUsage = dieTypeMap

	json.NewEncoder(w).Encode(response)
}

func diceRollHandler(w http.ResponseWriter, r *http.Request) {
	totalRequestCount++
	var rollList DieRollRequest
	err := json.NewDecoder(r.Body).Decode(&rollList)
	if err != nil {
		badRequestCount++
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

			count, ok := dieTypeMap[dieRoll.DieType]
			if ok {
				dieTypeMap[dieRoll.DieType] = count + dieRoll.RollCount
			} else {
				dieTypeMap[dieRoll.DieType] = dieRoll.RollCount
			}
		}
		goodRequestCount++
		json.NewEncoder(w).Encode(resultingRoll)
		return
	}

	badRequestCount++
	fmt.Fprintf(w, "No die rolls specs found in request")
}
