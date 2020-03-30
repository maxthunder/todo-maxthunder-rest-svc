package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
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

func getAllTasks(db *sql.DB) Tasks {
	return getFilteredTasks(db, true, true)
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

func getAllActiveTasks(db *sql.DB) Tasks {
	return getFilteredTasks(db, true, false)
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

func getAllCompletedTasks(db *sql.DB) Tasks {
	return getFilteredTasks(db, false, true)
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
	//json.NewEncoder(w).Encode(addNewTask(getDatabaseConnection(), task.Description))
	addNewTask(getDatabaseConnection(), task.Description)
}

func addNewTask(db *sql.DB, description string)  {
	taskSql := "INSERT INTO task(description, timestamp, isCompleted) VALUES(?, ?, false)"
	//taskSql := fmt.Sprintf("INSERT INTO task(description, timestamp, isCompleted) VALUES(%s, %s, false)", description, time.Now().String())
	rows, err := db.Query(taskSql, description, time.Now().String())
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	defer db.Close()


	//row := stmt.QueryRowContext(
	//	context.Background(),
	//	sql.Named("Description", description),
	//	sql.Named("Timestamp", time.Now().String()),
	//	sql.Named("IsCompleted", false))

	//stmt.Query(description, time.Now().String(), false)
}

//func addNewTask(db *sql.DB, description string)  {
//	taskSql := "INSERT INTO task(description, timestamp, isCompleted) VALUES(@Description, @Timestamp, @IsCompleted);"
//	stmt, err := db.Prepare(taskSql)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//
//	row := stmt.QueryRowContext(
//		context.Background(),
//		sql.Named("Description", description),
//		sql.Named("Timestamp", time.Now().String()),
//		sql.Named("IsCompleted", false))
//
//	fmt.Println(row)
//
//	stmt.Query(description, time.Now().String(), false)
//}

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
		var status string

		if includeActive && !task.IsCompleted {
			tasks = append(tasks, task)
		}

		if includeCompleted && task.IsCompleted {
			tasks = append(tasks, task)
		}


		fmt.Println(task.Description+ " (" + task.Timestamp + ") -> " + status)
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

func handleRequests() {
	Router := mux.NewRouter().StrictSlash(true)
	Router.HandleFunc("/", indexHandler)
	Router.HandleFunc("/tasks", getTasksHandler).Methods("GET", "OPTIONS")
	Router.HandleFunc("/activeTasks", postActiveTask).Methods("POST", "OPTIONS")
	Router.HandleFunc("/activeTasks", activeTasksHandler).Methods("GET")
	Router.HandleFunc("/completedTasks", completedTasksHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", Router))
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
