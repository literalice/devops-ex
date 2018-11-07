package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/gorilla/mux"

	"flag"
)

var resource string
var scanned = mapset.NewSet()
var winnersNumber = 5

func main() {
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)

	resource = "winners"
	root := "/api/" + resource

	router.HandleFunc("/", index)
	router.HandleFunc(root, winnersIndex)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "{\"message\": \"Welcome to the %s app!\" }", resource)
}

type result struct {
	Winners []string `json:"winners"`
}

func winnersIndex(w http.ResponseWriter, r *http.Request) {
	var candidates = strings.Split(os.Getenv("CANDIDATES"), ",")
	winnersNumber, _ := strconv.Atoi(r.URL.Query().Get("number"))
	winners := selectWinners(candidates, winnersNumber)

	result := result{
		Winners: winners,
	}

	jsonBytes, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	fmt.Fprintf(w, string(jsonBytes))
}

func selectWinners(candidates []string, maxWinnerNumber int) []string {
	candidatesNumber := len(candidates)
	if candidatesNumber <= maxWinnerNumber {
		return candidates
	}

	rand.Seed(time.Now().UnixNano())
	winners := make([]string, 0, maxWinnerNumber)
	for len(winners) < maxWinnerNumber && len(candidates) > 0 {
		candidateIndex := rand.Intn(len(candidates))
		candidate := candidates[candidateIndex]

		if !scanned.Contains(candidate) {
			winners = append(winners, candidate)
		}
		scanned.Add(candidate)

		candidates = unset(candidates, candidateIndex)
	}
	return winners
}

func unset(s []string, i int) []string {
	if i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}
