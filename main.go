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
	TaskId string `json:"taskID"`
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
	tasks := getAllActiveTasks(getDatabaseConnection())

	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	fmt.Println("("+time.Now().String()+") Endpoint Hit: /activeTasks")
	json.NewEncoder(w).Encode(tasks)
}

func getAllCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := getAllCompletedTasks(getDatabaseConnection())

	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	fmt.Println("("+time.Now().String()+") Endpoint Hit: /completedTasks")
	json.NewEncoder(w).Encode(tasks)
}

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", "todoDatasource_user:todoDatasource_user123@tcp(127.0.0.1:3306)/todoDatasource")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getAllTasks(db *sql.DB) Tasks {
	return getTasks(db, true, true)
}

func getAllActiveTasks(db *sql.DB) Tasks {
	return getTasks(db, true, false)
}

func getAllCompletedTasks(db *sql.DB) Tasks {
	return getTasks(db, false, true)
}

func getTasks(db *sql.DB, includeActive bool, includeCompleted bool) Tasks {
	results, err := db.Query("SELECT * FROM task")
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
		var status string

		if includeActive && !task.IsCompleted {
			tasks = append(tasks, task)
		}

		if includeCompleted && task.IsCompleted {
			tasks = append(tasks, task)
		}


		fmt.Println(task.Description+ " (" + task.Timestamp + ") -> " + status)
	}
	fmt.Printf("Number of tasks : %v\n", len(tasks))
	defer db.Close()
	return tasks
}

func handleRequests() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/activeTasks", getAllActiveTasksHandler)
	http.HandleFunc("/completedTasks", getAllCompletedTasksHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func insertTestData(db *sql.DB) {
	insert, err := db.Query("INSERT INTO task (description, timestamp, isCompleted) VALUES ('Test task', now(), false)")
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	fmt.Println("Successful INSERT into 'task' table.")

}

func main() {
		handleRequests()
}
