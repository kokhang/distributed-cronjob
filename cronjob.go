package main

import (
	"encoding/json"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/kokhang/distributed-cronjob/manager"
	"github.com/kokhang/distributed-cronjob/models"
)

func GetJobEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range models.Jobs {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&models.Job{})
}

func ListJobsEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(models.Jobs)
}

func CreateJobEndpoint(w http.ResponseWriter, req *http.Request) {
	var job models.Job
	_ = json.NewDecoder(req.Body).Decode(&job)

	// TODO validate(job)
	manager.ScheduleNewJob(job)
	json.NewEncoder(w).Encode(models.Jobs)
}

func main() {
	router := mux.NewRouter()
	manager.InitWorkers()
	router.HandleFunc("/jobs", ListJobsEndpoint).Methods("GET")
	router.HandleFunc("/jobs/{id}", GetJobEndpoint).Methods("GET")
	router.HandleFunc("/jobs", CreateJobEndpoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":30000", router))
}
