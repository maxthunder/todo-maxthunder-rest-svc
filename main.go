package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"os"
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

// GET /status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	json.NewEncoder(w).Encode("200 OK")
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

// POST /tasks
func postActiveTask(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: POST /tasks")


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

// PUT /tasks
func updateActiveTask(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: PUT /tasks")

	decoder := json.NewDecoder(r.Body)
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		panic(err.Error())
	}
	if updateTask(getDatabaseConnection(), task) {
		json.NewEncoder(w).Encode("Task with ID " + string(task.TaskId) + " was successfully completed.")
	}
}

// DELETE /tasks
func deleteTask(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Println("("+time.Now().String()+") Endpoint Hit: DELETE /tasks")

	taskId, ok := r.URL.Query()["taskId"]
	if !ok || len(taskId[0]) < 1 {
		panic("Url request parameter 'taskId' is required for task deletion.")
	}
	if deleteCompletedTask(getDatabaseConnection(), taskId[0]) {
		json.NewEncoder(w).Encode("Task with ID " + taskId[0] + " was successfully deleted.")
	}
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


// Database Utilities
func addNewTask(db *sql.DB, description string) bool  {
	taskSql := "INSERT INTO todo.task(taskId, description, timestamp, iscompleted) VALUES(default, $1, $2, false);"

	var now = time.Now()
	var taskDate = fmt.Sprintf("%v %v %v", now.Month().String(), now.Day(), now.Year())
	results, err := db.Query(taskSql, description, taskDate)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	defer db.Close()
	return true
}

func completeTask(db *sql.DB, taskId int) bool {
	results, err := db.Query("UPDATE todo.task SET iscompleted=true WHERE taskId=$1", taskId)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	defer db.Close()
	return true
}

func updateTask(db *sql.DB, task Task) bool {
	results, err := db.Query("UPDATE todo.task SET description=$1, timestamp=$2, iscompleted=$3 WHERE taskId=$4",
		task.Description, task.Timestamp, task.IsCompleted, task.TaskId)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	defer db.Close()
	return true
}

func deleteCompletedTask(db *sql.DB, taskId string) bool {
	results, err := db.Query("DELETE FROM todo.task WHERE taskId=$1", taskId)
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
		query = "SELECT * FROM todo.task"
	} else if includeActive {
		query = "SELECT * FROM todo.task WHERE iscompleted = false"
	} else if includeCompleted {
		query = "SELECT * FROM todo.task WHERE iscompleted = true"
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
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "postgres://prlgrktludlpfy:1723a94575704248af1d99b8683452ee3de33b48d88b18856c4109662b41b995@ec2-34-206-252-187.compute-1.amazonaws.com:5432/d3tvcfn8ldd2av")
	//db, err := sql.Open("mysql", "todoDatasource_user:todoDatasource_user123@tcp(127.0.0.1:3306)/todoDatasource")
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
	Router.HandleFunc("/status", statusHandler).Methods("GET", "OPTIONS")
	Router.HandleFunc("/tasks", getTasksHandler).Methods("GET", "OPTIONS")
	Router.HandleFunc("/tasks", postActiveTask).Methods("POST", "OPTIONS")
	Router.HandleFunc("/tasks", updateActiveTask).Methods("PUT", "OPTIONS")
	Router.HandleFunc("/tasks", deleteTask).Methods("DELETE", "OPTIONS")
	Router.HandleFunc("/activeTasks", activeTasksHandler).Methods("GET", "OPTIONS")
	Router.HandleFunc("/completedTasks", completedTasksHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), Router))
}

func main() {
		handleRequests()
}
