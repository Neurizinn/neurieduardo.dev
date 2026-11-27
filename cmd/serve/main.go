package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Serving dist at http://localhost:8080; rebuild with make build after edits")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("dist"))))
}
