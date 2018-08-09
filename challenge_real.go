package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var get_count uint64
var create_count uint64
var load_count uint64

type Card struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Load struct {
	ReferenceID string  `json:"reference_id"`
	Amount      float32 `json:"amount"`
}

type Response struct {
	Data struct {
		NameOnCard  string    `json:"name_on_card"`
		Pan         string    `json:"pan"`
		ReferenceID string    `json:"reference_id"`
		ExpDate     string    `json:"exp_date"`
		Balance     int       `json:"balance"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"data"`
}

type ResponseError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type Fail struct {
	Message []byte
	Intent  int
}

func Newfail(message []byte) Fail {

	fail := Fail{}
	fail.Message = message
	fail.Intent = 0
	return fail
}

func Fooget() {
	atomic.AddUint64(&get_count, 1)
	fmt.Println("+1!!!!!")
}
func Foocreate() {
	atomic.AddUint64(&create_count, 1)
	fmt.Println("+1!!!!!")
}
func Fooload() {
	atomic.AddUint64(&load_count, 1)
}

var create = make(chan Fail, 100)
var load = make(chan Fail, 100)
var get = make(chan Fail, 100)

var getchan = make(chan string, 2)

func createCard(w http.ResponseWriter, r *http.Request) {

	url := "https://fakeprovider.herokuapp.com/cards"

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fail := Newfail(body)

	create <- fail
	Foocreate()
	fmt.Println(create_count)
	if create_count%2 == 0 && create_count != 0 {
		fmt.Println("tuto wawa")
		time.Sleep(time.Second * 10)
		fmt.Println("Desperto")
	}
	go sendjson(create, url, "card")

	defer r.Body.Close()

}

func loadCard(w http.ResponseWriter, r *http.Request) {

	url := "https://fakeprovider.herokuapp.com/load"

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fail := Newfail(body)

	load <- fail
	Fooload()
	fmt.Println(load_count)
	if load_count%2 == 0 && load_count != 0 {
		fmt.Println("tuto wawa")
		time.Sleep(time.Second * 10)
		fmt.Println("Desperto")
	}
	go sendjson(load, url, "load")

	defer r.Body.Close()

}

func sendjson(c chan Fail, url string, action string) {

	d := <-c

	if action == "card" {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(d.Message))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println("response Status:", resp.Status)
		errorHandler(resp.Status, d, action)
		fmt.Printf("%#v\n", resp)

	} else if action == "get" {

		Fooget()
		fmt.Println(get_count)

		if get_count%2 == 0 && get_count != 0 {
			fmt.Println("tuto wawa")
			time.Sleep(time.Second * 10)
			fmt.Println("Desperto")
		}

		req, err := http.Get(url)
		if err != nil {
			log.Print(err.Error())
			os.Exit(1)
		}

		fmt.Println("response Status:", req.Status)
		fmt.Println("salio")
		errorHandler(req.Status, <-c, action)
		responseData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(responseData))
	} else if action == "load" {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(d.Message))
		req.Header.Set("Authorization", "Bearer fasdfadfa9fj987afsdf")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println("response Status:", resp.Status)
		errorHandler(resp.Status, d, action)
		fmt.Printf("%#v\n", resp)
	}

}

func getCard(w http.ResponseWriter, r *http.Request) {

	url := "https://fakeprovider.herokuapp.com/"

	fail := Newfail(nil)

	get <- fail
	fmt.Println("entro a la cola")

	go sendjson(get, url, "get")

}

func errorHandler(err string, p Fail, action string) {
	fmt.Println("que sale en el error")
	fmt.Println(p.Intent)
	if p.Intent >= 5 {
		fmt.Println("a fallado muchas veces")

	} else if action == "get" {
		switch {
		case err == "500 Internal Server Error":
			p.Intent++
			get <- p
			go sendjson(get, "https://fakeprovider.herokuapp.com/", "get")

		case err == "400 Bad Request":

		case err == "404 Not Found":
			p.Intent++
			get <- p
			go sendjson(get, "https://fakeprovider.herokuapp.com/", "get")

		case err == "401 Unauthorized":

		case err == "429 Too Many Requests":

			fmt.Println("SE ENCONTRO UN ERROROROROROROR")
			p.Intent++
			get <- p
			go sendjson(get, "https://fakeprovider.herokuapp.com/", "get")

		}
	} else if action == "card" {
		switch {
		case err == "500 Internal Server Error":
			p.Intent++
			create <- p

		case err == "400 Bad Request":

		case err == "404 Not Found":
			p.Intent++
			create <- p
		case err == "401 Unauthorized":

		case err == "429 Too Many Requests":
			p.Intent++
			create <- p

		}
	} else if action == "load" {
		switch {
		case err == "500 Internal Server Error":
			p.Intent++
			create <- p

		case err == "400 Bad Request":

		case err == "404 Not Found":
			p.Intent++
			create <- p
		case err == "401 Unauthorized":

		case err == "429 Too Many Requests":
			p.Intent++
			create <- p

		}
	}

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", getCard).Methods("GET")
	router.HandleFunc("/cards", createCard).Methods("POST")
	router.HandleFunc("/load", loadCard).Methods("POST")
	//	router.HandleFunc("/cards/:id/info").Methods("")

	log.Fatal(http.ListenAndServe(":8080", router))

}
