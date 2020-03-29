package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type Task struct {
	Description string `json:"description"`
	//Timestamp time.Time `json:"timestamp"`
	Timestamp string `json:"timestamp"`
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

func getAllActiveTasksHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	tasks := Tasks{
		Task{Description:"*Feed the fish.", Timestamp:"03/15/2020", IsCompleted:false},
		Task{Description:"*Vacuum the carpet.", Timestamp:"03/20/2020", IsCompleted:false},
	}

	fmt.Println("("+time.Now().String()+") Endpoint Hit: /activeTasks")
	json.NewEncoder(w).Encode(tasks)
}

func getAllCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	tasks := Tasks{
		Task{Description:"*Dust the table.", Timestamp:"03/20/2020", IsCompleted:true},
		Task{Description:"*Clean the sink.", Timestamp:"03/31/2020", IsCompleted:true},
	}

	fmt.Println("("+time.Now().String()+") Endpoint Hit: /completedTasks")
	json.NewEncoder(w).Encode(tasks)
}

func handleRequests() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/activeTasks", getAllActiveTasksHandler)
	http.HandleFunc("/completedTasks", getAllCompletedTasksHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	db, err := sql.Open("mysql", "todoDatasource_user:todoDatasource_user123@tcp(127.0.0.1:3306)/todoDatasource")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully connected to MySQL database.")

	handleRequests()
}
