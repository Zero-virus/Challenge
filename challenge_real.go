package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	var c = make(chan io.ReadCloser, 100)
	url := "https://fakeprovider.herokuapp.com/cards"

	go createqueue(r, c)
	// decoder := json.NewDecoder(r.Body)
	// var t Card

	fmt.Println("llego aqui")

	// err := decoder.Decode(&t)
	// if err != nil {

	// }
	//json.NewEncoder(w).Encode(t)
	// jsn, err := ioutil.ReadAll(r.Body)

	resp, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		fmt.Println("error", err)
		return
	} //else {
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	//}
	//var resperr ResponseError
	//json.Unmarshal(resp, &resperr)

	fmt.Printf("%#v\n", resp)
	sendjson(c, w, r, url)
}

func createqueue(r *http.Request, c chan io.ReadCloser) {

	c <- r.Body
	//count := 0

	close(c)

}

func sendjson(c chan io.ReadCloser, w http.ResponseWriter, r *http.Request, url string) {

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
