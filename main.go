package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Simple service to test how Sesam behaves when sending, transforming or getting data under different  circumstances
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Starting service on port %s", port)

	router := mux.NewRouter()
	router.HandleFunc("/get", FetchData).Methods("GET")
	router.HandleFunc("/post", SendData).Methods("POST")
	router.HandleFunc("/transform", TransformData).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func FetchData(w http.ResponseWriter, r *http.Request) {
//TODO
}

//Send data simulates data sink
//This endpoint takes array of JSON elements on input and return nothing (or error)
//input data available options:
//	*should_fail - if exists and equal to true then 400 bad request will be returned back
func SendData(w http.ResponseWriter, r *http.Request) {
	log.Println("\r\n====\t\tREQUEST START\t\t====")
	var bodyJsonArray []interface{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(body, &bodyJsonArray)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Got batch of %d entities", len(bodyJsonArray))

	for _, item := range bodyJsonArray {
		mappedItem := item.(map[string]interface{})
		log.Printf("Processing entity with id: %s", mappedItem["_id"])
		//will return 400 Bad request if entity contains this key evaluated as true
		if mappedItem["should_fail"] != nil && mappedItem["should_fail"] == true {
			log.Printf("Entity %s should fail, return HTTP 400", mappedItem["_id"])
			http.Error(w, "Entity had to fail =(", http.StatusBadRequest)
			break
		}
	}

	log.Println("\r\n====\t\tREQUEST FINISH\t\t====")
}

func TransformData(w http.ResponseWriter, r *http.Request) {
//TODO
}
