package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Test struct {
	ID   int    `json:"id"`
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
	json.NewEncoder(w).Encode(<-c)
	fmt.Println("salio")
}

func handler(w http.ResponseWriter, r *http.Request) {

	var c = make(chan Test, 2)

	t := Test{
		Name: "bab",
	}

	count := 0
	go func() {
		for n := 0; n < 5; n++ {
			count++
			t.ID = n
			c <- t
			fmt.Println("entro", n)

		}
		close(c)
	}()
	fmt.Println(len(c))

	enviad := 0
	for d := 0; d < 5; d++ {
		sendjson(c, w, r)
		enviad++
		if enviad%2 == 0 {
			fmt.Println("tuto wawa")
			time.Sleep(time.Second * 10)
			fmt.Println("Desperto")
		}
	}
	/*
		go func() {
			c <- t
			close(c)
		}()

		}*/

}

func receive(w http.ResponseWriter, r *http.Request) {

}

func main() {

	router := mux.NewRouter()

	//	c := gener()
	fmt.Println("cerrado")

	router.HandleFunc("/", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}
