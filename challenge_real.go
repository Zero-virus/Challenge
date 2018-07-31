package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Card struct {
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
}

type Load struct {
	Reference_id string  `json:"reference_id"`
	Amount       float32 `json:"amount"`
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

func createCard(w http.ResponseWriter, r *http.Request) {

	createqueue(r)
}

func createqueue(r *http.Request) {
	//var c = make(chan http.ResponseWriter, 2)
	//	url := "https://fakeprovider.herokuapp.com"

}

func getCard(w http.ResponseWriter, r *http.Request) {
	url := "https://fakeprovider.herokuapp.com/"
	req, err := http.Get(url)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
	var responseObj Response
	json.Unmarshal(responseData, &responseObj)
	fmt.Println(responseObj)
}

func loadCard(w http.ResponseWriter, r *http.Request) {

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", getCard).Methods("GET")
	router.HandleFunc("/cards", createCard).Methods("POST")
	router.HandleFunc("/load", loadCard).Methods("POST")
	//	router.HandleFunc("/cards/:id/info").Methods("")

	log.Fatal(http.ListenAndServe(":8080", router))

}
