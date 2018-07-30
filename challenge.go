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

/*func gener(te []Test) <-chan Test {
	out := make(chan []Test)
	go func() {
		for _, n := range te {
			out <- te
		}
		close(out)
	}()
	fmt.Println("llamar return")
	return out
}*/

func sendjson(c chan Test, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	v := <-c
	json.NewEncoder(w).Encode(v)
}

func handler(w http.ResponseWriter, r *http.Request) {

	//	var t Test
	//	json.NewDecoder(r.Body).Decode(&t)

	var c = make(chan Test)

	t := []Test{
		{
			ID:   "daad",
			Name: "bab",
		},
		{
			ID:   "daa",
			Name: "bb",
		},
	}

	count := 0
	go func() {
		for _, n := range t {
			count++
			c <- n

		}
		close(c)
	}()
	fmt.Println(len(c))

	for d := 0; d < count; d++ {
		go sendjson(c, w, r)
	}
	/*
		go func() {
			c <- t
			close(c)
		}()

		}*/

}

func receive()

func main() {

	router := mux.NewRouter()

	//	c := gener()
	fmt.Println("cerrado")

	router.HandleFunc("/", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}
