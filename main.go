package main

import (
	"net/http"
	"log"
)

func main() {
	fs := http.FileServer(http.Dir("www"))
	log.Fatal(http.ListenAndServe(":8080", fs))
}
