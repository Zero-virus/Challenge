package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

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

func createCard(w http.ResponseWriter, r *http.Request) {

	url := "https://fakeprovider.herokuapp.com/cards"

	//n := r.Body
	//	c <- n
	count := 0
	var card Card

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println(card)
	go createqueue(body, url, count, "card")

	defer r.Body.Close()

}

func loadCard(w http.ResponseWriter, r *http.Request) {

	url := "https://fakeprovider.herokuapp.com/load"

	count := 0

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	go createqueue(body, url, count, "load")

	defer r.Body.Close()

}

func createqueue(b []byte, url string, count int, action string) {
	var c = make(chan []byte, 100)
	//count := 0
	c <- b
	close(c)

	fmt.Println(count)

	if count%2 == 0 {
		fmt.Println("tuto wawa")
		time.Sleep(time.Second * 10)
		fmt.Println("Desperto")
	}
	sendjson(c, url, action)
}

func sendjson(c chan []byte, url string, action string) {
	d := <-c

	if action == "card" {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println("response Status:", resp.Status)
		errorHandler(resp.Status, d)
		fmt.Printf("%#v\n", resp)

	} else if action == "get" {
		req, err := http.Get(url)
		if err != nil {
			log.Print(err.Error())
			os.Exit(1)
		}

		fmt.Println("response Status:", req.Status)
		fmt.Println("salio")
		responseData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(responseData))
	} else if action == "load" {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
		req.Header.Set("Authorization", "Bearer fasdfadfa9fj987afsdf")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println("response Status:", resp.Status)
		errorHandler(resp.Status, d)
		fmt.Printf("%#v\n", resp)
	}

}

func getCard(w http.ResponseWriter, r *http.Request) {
	url := "https://fakeprovider.herokuapp.com/"
	count := 0
	go createqueue(nil, url, count, "get")

}

func errorHandler(err string, p []byte) {

	var e chan []byte

	switch {
	case err == "500 Internal Server Error":
		e <- p

	case err == "400 Bad Request":

	case err == "404 Not Found":
		e <- p
	case err == "401 Unauthorized":

	case err == "429 Too Many Request":
		e <- p

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
