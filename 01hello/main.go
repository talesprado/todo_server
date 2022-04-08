package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Todo struct {
	TodoId      int    `json:todoId`
	Text        string `json:"text"`
	IsCompleted bool   `json:"isCompleted"`
}
type Todos []Todo

var todos = Todos{
	Todo{TodoId: 1, Text: "Fazer compras", IsCompleted: false},
	Todo{TodoId: 2, Text: "Lavar carro", IsCompleted: false},
	Todo{TodoId: 3, Text: "Estudar", IsCompleted: false},
	Todo{TodoId: 4, Text: "Fazer compras", IsCompleted: false},
	Todo{TodoId: 5, Text: "Lavar carro", IsCompleted: false},
	Todo{TodoId: 6, Text: "Estudar", IsCompleted: false},
}

func allTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId, _ := strconv.Atoi(vars["id"])

	for key, value := range todos {
		if value.TodoId == todoId {
			json.NewEncoder(w).Encode(todos[key])
		}
	}
}

func completeTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId, _ := strconv.Atoi(vars["id"])

	for key, value := range todos {
		if value.TodoId == todoId {
			todos[key].IsCompleted = true
			json.NewEncoder(w).Encode(todos[key])
		}
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId, _ := strconv.Atoi(vars["id"])

	for key, value := range todos {
		if value.TodoId == todoId {
			todos = append(todos[:key], todos[key+1:]...)
		}
	}
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			var newId = 1
			if len(todos) > 0 {
				var last = todos[len(todos)-1]
				newId = last.TodoId + 1
			}

			todo.TodoId = newId
			todo.IsCompleted = false
			todos = append(todos, todo)
			json.NewEncoder(w).Encode(todos[len(todos)-1])
		}
	}

}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})

	myRouter.HandleFunc("/todos", allTodos).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/todo/{id}", getTodo).Methods("GET")
	myRouter.HandleFunc("/todo/delete/{id}", deleteTodo).Methods("GET")
	myRouter.HandleFunc("/todo/complete/{id}", completeTodo).Methods("GET")
	myRouter.HandleFunc("/todo/add", addTodo).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(headers, methods, origins)(myRouter)))
}
