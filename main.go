package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Congratulations! Welcome to the Swisscom Application Cloud!")
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
