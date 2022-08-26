package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var greeting string

func handler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Handling request: [%s; %s; %s]\n", req.Method, req.RequestURI, req.UserAgent())

	if req.URL.Path == "/" {
		fmt.Fprintf(rw, "%s!", greeting)
	} else {
		fmt.Fprintf(rw, "%s, %s!", greeting, req.URL.Path[1:])
	}
}

func main() {
	// can be set by manifest.yml during "cf push", or by "cf set-env"
	greeting = os.Getenv("GREETING")
	if len(greeting) == 0 {
		greeting = "Hello"
	}

	// read VCAP_APP_HOST and PORT env variables set by CloudFoundry
	host := os.Getenv("VCAP_APP_HOST")
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	address := host + ":" + port

	http.HandleFunc("/", handler)

	log.Printf("Starting web server, listening on [%s] ...\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}
