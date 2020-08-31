package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RawAlert struct {
	RuleName string
	Tags     Tags
	Title    string
	Message  string
	State    string
}
type Tags struct {
	Tag      string
	Priority string
}

func parseAlert(alert string) RawAlert {
	var alertObj RawAlert
	if err := json.Unmarshal([]byte(alert), &alertObj); err != nil {
		panic(err)
	}
	return alertObj
}

func status(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "up and running!")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
			log.Printf("%v: %v\n", name, h)
		}
		log.Println(req.Body)
		/*
			decoder := json.NewDecoder(req.Body)
			var t test_struct
			err := decoder.Decode(&t)
			if err != nil {
				panic(err)
			}
			log.Println(t.Test)
		*/
	}
}

func main() {
	addr := ":8090"
	http.HandleFunc("/status", status)
	//http.HandleFunc("/headers", headers)
	log.Println("listen on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
