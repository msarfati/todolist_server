package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID   int    `json:"id,omitempty"`
	Text string `json:"text"`
}

var todoList []Todo

func NewIndex(t []Todo) int {
	max := 0
	for _, item := range t {
		if int(item.ID) > max {
			max = item.ID
		}
	}
	return max + 1
}

func GetTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	paramsID, _ := strconv.Atoi(params["id"])
	for _, item := range todoList {
		if item.ID == paramsID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Todo{})
}

func GetTodoListEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(todoList)
}

func CreateTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	var todo Todo
	_ = json.NewDecoder(req.Body).Decode(&todo)
	todo.ID = NewIndex(todoList)
	todoList = append(todoList, todo)
	json.NewEncoder(w).Encode(todoList)
}

func DeleteTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	paramsID, _ := strconv.Atoi(params["id"])
	for index, item := range todoList {
		if item.ID == paramsID {
			todoList = append(todoList[:index], todoList[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(todoList)
}

func main() {
	router := mux.NewRouter()
	todoList = append(todoList, Todo{ID: 1, Text: "Do laundry."})
	router.HandleFunc("/todo", GetTodoListEndpoint).Methods("GET")
	router.HandleFunc("/todo/{id}", GetTodoEndpoint).Methods("GET")
	router.HandleFunc("/todo", CreateTodoEndpoint).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteTodoEndpoint).Methods("DELETE")
	fmt.Println("Server running on port 12345")
	log.Fatal(http.ListenAndServe(":12345", router))
}
