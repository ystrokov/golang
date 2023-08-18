package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Person struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Developer string `json:"developer"`
}

var (
	people     []Person
	idSequence int
	mutex      sync.Mutex
)

func addPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Автоматическое присвоение нового id
	idSequence++
	person.ID = idSequence

	if person.Name == "" || person.Surname == "" {
		http.Error(w, "Name and surname are required", http.StatusBadRequest)
		return
	}

	for _, p := range people {
		if p.Name == person.Name && p.Surname == person.Surname {
			http.Error(w, "Person already exists", http.StatusBadRequest)
			return
		}
	}

	people = append(people, person)
	fmt.Fprintf(w, "Person added: %s %s\n", person.Name, person.Surname)
	w.WriteHeader(http.StatusCreated)
}

func getPeopleHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func deletePersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var index int = -1
	for i, p := range people {
		if fmt.Sprintf("%d", p.ID) == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Person with specified id not found", http.StatusNotFound)
		return
	}

	people = append(people[:index], people[index+1:]...)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/add", addPersonHandler)
	http.HandleFunc("/get", getPeopleHandler)
	http.HandleFunc("/delete", deletePersonHandler)

	fmt.Println("Server i`s running on :8080")
	http.ListenAndServe(":8080", nil)
}
