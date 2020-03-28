package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Task struct {
	Description string `json:"description"`
	Timestamp time.Time `json:"timestamp"`
	IsCompleted bool `json:"isCompleted"`
}

type Tasks []Task

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	tasks := Tasks{
		Task{Description:"Feed the fish.", Timestamp:time.Now(), IsCompleted:false},
		Task{Description:"Vacuum the carpet.", Timestamp:time.Now(), IsCompleted:false},
		Task{Description:"Clean the sink.", Timestamp:time.Now(), IsCompleted:true},
	}

	fmt.Println("Endpoint Hit: getTasks()")
	json.NewEncoder(w).Encode(tasks)
}

func handleRequests() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tasks", getAllTasksHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
