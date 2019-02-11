package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Simple service to test how Sesam behaves when sending, transforming or getting data under different  circumstances
func main() {
	port := os.Getenv("PORT")
	rand.Seed(time.Now().UnixNano())

	if port == "" {
		port = "8080"
	}

	log.Printf("Starting service on port %s", port)

	router := mux.NewRouter()
	router.HandleFunc("/get/{howMany}", FetchData).Methods("GET")
	router.HandleFunc("/post", SendData).Methods("POST")
	router.HandleFunc("/transform", TransformData).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func FetchData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving request %v\r\n", r.URL)
	vars := mux.Vars(r)
	howMany, err := strconv.Atoi(vars["howMany"])
	if err != nil {
		howMany = 0
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))
	first := true

	for i := 0; i < howMany; i++ {
		if first {
			first = false
		} else {
			w.Write([]byte(","))
		}
		item := GenerateItem()
		log.Printf("Generated item %v", item)
		jsonData, err := json.Marshal(item)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	}
	w.Write([]byte("]"))
}

type Item struct {
	Id         string `json:"_id"`
	ShouldFail bool   `json:"should_fail"`
	Value      string `json:"value"`
}

func GenerateItem() Item {
	inputDataArr := []string{"foo", "bar", "baz", "karabas"}
	var result Item
	result.Id = inputDataArr[rand.Intn(len(inputDataArr))] + strconv.Itoa(rand.Intn(9999999))
	result.Value = inputDataArr[rand.Intn(len(inputDataArr))]

	if rand.Intn(10) > 8 {
		result.ShouldFail = true
	} else {
		result.ShouldFail = false
	}

	return result
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
