package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gihub.com/meli-dna/dna"
	"gihub.com/meli-dna/repository"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type request struct {
	Dna []string `json:"dna"`
}

type statsResponse struct {
	CountMutantDNA int64   `json:"count_mutant_dna"`
	CountHumanDNA  int64   `json:"count_human_dna"`
	Ratio          float32 `json:"ratio"`
}

var repo repository.Interface
var dnaService dna.Interface
var server http.Server

var ctx = context.Background()

func init() {
	dnaService = dna.Service{}
	repo = repository.Repository{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})}
}

func isMutant(w http.ResponseWriter, r *http.Request) {
	var newRequest request

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter a dna array")
	}

	json.Unmarshal(reqBody, &newRequest)

	isMutant := dnaService.IsMutant(newRequest.Dna)

	repo.SetDNA(newRequest.Dna, isMutant)

	if isMutant {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func stats(w http.ResponseWriter, r *http.Request) {
	mutantCount, humanCount := repo.GetDNACounts()

	ratio := float32(mutantCount) / float32(mutantCount+humanCount)

	statsResponse := statsResponse{
		CountMutantDNA: mutantCount,
		CountHumanDNA:  humanCount,
		Ratio:          ratio,
	}

	json.NewEncoder(w).Encode(statsResponse)
}

// Start initialize webserver
func Start() {
	repo.Initialize()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mutant", isMutant).Methods("POST")
	router.HandleFunc("/stats", stats).Methods("GET")

	server.Addr = ":8080"
	server.Handler = router

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
