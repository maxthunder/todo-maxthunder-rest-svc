package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Task struct {
	TaskId int `json:"taskId"`
	Description string `json:"description"`
	Timestamp string `json:"timestamp"`
	IsCompleted bool `json:"isCompleted"`
}

type Tasks []Task


// GET /
func indexHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

// GET /tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: GET /tasks")

	tasks := getAllTasks(getDatabaseConnection())

	json.NewEncoder(w).Encode(tasks)
}

// GET /activeTasks
func activeTasksHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: GET /activeTasks")

	var tasks = getAllActiveTasks(getDatabaseConnection())

	json.NewEncoder(w).Encode(tasks)
}

// GET /completedTasks
func completedTasksHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: GET /completedTasks")

	tasks := getAllCompletedTasks(getDatabaseConnection())

	json.NewEncoder(w).Encode(tasks)
}

// POST /activeTasks
func postActiveTask(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: POST /activeTasks")


	//tasks := addNewTask(getDatabaseConnection())
	decoder := json.NewDecoder(r.Body)
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		panic(err.Error())
	}
	if addNewTask(getDatabaseConnection(), task.Description) {
		json.NewEncoder(w).Encode("Task with description " + task.Description + " was successfully added.")
	}
}

// PUT /activeTasks
func updateActiveTask(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: PUT /activeTasks")

	decoder := json.NewDecoder(r.Body)
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		panic(err.Error())
	}
	if completeTask(getDatabaseConnection(), task.TaskId) {
		json.NewEncoder(w).Encode("Task with ID " + string(task.TaskId) + " successfully completed.")
	}
}

// Database Utilities
func addNewTask(db *sql.DB, description string) bool  {
	taskSql := "INSERT INTO task(description, timestamp, isCompleted) VALUES(?, ?, false)"
	var now = time.Now()
	var taskDate = fmt.Sprintf("%v %v, %v", now.Month().String(), now.Day(), now.Year())
	results, err := db.Query(taskSql, description, taskDate)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	defer db.Close()
	return true
}

func completeTask(db *sql.DB, taskId int) bool {
	results, err := db.Query("UPDATE task SET isCompleted=true WHERE taskId=?", taskId)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	defer db.Close()
	return true
}

func getAllTasks(db *sql.DB) Tasks {
	return getFilteredTasks(db, true, true)
}

func getAllActiveTasks(db *sql.DB) Tasks {
	return getFilteredTasks(db, true, false)
}

func getAllCompletedTasks(db *sql.DB) Tasks {
	return getFilteredTasks(db, false, true)
}

func getFilteredTasks(db *sql.DB, includeActive bool, includeCompleted bool) Tasks {
	var query string

	if includeActive && includeCompleted {
		query = "SELECT * FROM task"
	} else if includeActive {
		query = "SELECT * FROM task WHERE isCompleted = false"
	} else if includeCompleted {
		query = "SELECT * FROM task WHERE isCompleted = true"
	} else {
		panic("includeActive & includeActive cannot both be false")
	}

	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	tasks := Tasks{}

	for results.Next() {
		var task Task

		err = results.Scan(&task.TaskId,&task.Description,&task.Timestamp,&task.IsCompleted)
		if err != nil {
			panic(err.Error())
		}

		// append on active tasks for active task searches OR append on completed task for completed task searches.
		if (includeActive && !task.IsCompleted) || (includeCompleted && task.IsCompleted) {
			tasks = append(tasks, task)
		}

		fmt.Println("(" + task.Timestamp + ") : " + task.Description)
	}
	fmt.Printf("Number of tasks returned: %v\n", len(tasks))
	defer db.Close()
	return tasks
}

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", "todoDatasource_user:todoDatasource_user123@tcp(127.0.0.1:3306)/todoDatasource")
	//db, err := sql.Open("sqlserver", "todoDatasource_user:todoDatasource_user123@tcp(127.0.0.1:3306)/todoDatasource")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Http and REST Utilities
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handleRequests() {
	Router := mux.NewRouter().StrictSlash(true)
	Router.HandleFunc("/", indexHandler)
	Router.HandleFunc("/tasks", getTasksHandler).Methods("GET", "OPTIONS")
	Router.HandleFunc("/activeTasks", postActiveTask).Methods("POST", "OPTIONS")
	Router.HandleFunc("/activeTasks", updateActiveTask).Methods("PUT", "OPTIONS")
	Router.HandleFunc("/activeTasks", activeTasksHandler).Methods("GET")
	Router.HandleFunc("/completedTasks", completedTasksHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", Router))
}

func main() {
		handleRequests()
}
