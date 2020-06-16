package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://nrk.no")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
