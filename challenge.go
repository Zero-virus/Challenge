package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Test struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func gener(te []Test) <-chan Test {
	out := make(chan []Test)
	go func() {
		for _, n := range te {
			out <- n
		}
		close(out)
	}()
	return out
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var t Test
	json.NewDecoder(r.Body).Decode(&t)

	var c = make(chan *Test)

	t := &Test{
		ID:   "daad",
		Name: "bab",
	}

	go func() {
		c <- t
		close(c)
	}()
	for n := range c {
		fmt.Println(w, string(n))
	}
}

func main() {

	router := mux.NewRouter()

	fmt.Println("cerrado")

	router.HandleFunc("/", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}
