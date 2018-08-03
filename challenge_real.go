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

	n := r.Body
	c <- n

	var card Card
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&card)
	if err != nil {
		panic(err)
	}

	fmt.Println(card)
	go createqueue(c, url)

	defer r.Body.Close()

	fmt.Println("llego aqui")

	// resp, err := http.Post(url, "application/json", r.Body)
	// if err != nil {
	// 	fmt.Println("error", err)
	// 	return
	// }
	// data, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(data))
	// fmt.Printf("%#v\n", resp)
}

func loadCard(w http.ResponseWriter, r *http.Request) {

	url := "https://fakeprovider.herokuapp.com/load"

	resp, err := http.Post(url, "", r.Body)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	fmt.Printf("%#v\n", resp)

}

func createqueue(c chan io.ReadCloser, url string) {

	close(c)

	sendjson(c, url)
}

func sendjson(c chan io.ReadCloser, url string) {
	d := <-c

	var card Card
	decoder := json.NewDecoder(d)
	err := decoder.Decode(&card)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(card)

	resp, err := http.Post(url, "application/json", d)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	fmt.Printf("%#v\n", resp)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseObj)
	fmt.Println("salio")

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", getCard).Methods("GET")
	router.HandleFunc("/cards", createCard).Methods("POST")
	router.HandleFunc("/load", loadCard).Methods("POST")
	//	router.HandleFunc("/cards/:id/info").Methods("")

	log.Fatal(http.ListenAndServe(":8080", router))

}
