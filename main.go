package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

// Kitten is a cute kitty
type Kitten struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name string        `json:"name"`
}

var dbString string

// kittenHandler allows to manage kittnes
func kittenHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(dbString)
	if err != nil {
		handleError(w, err)
	}
	defer session.Close()
	c := session.DB("").C("kittens")

	switch r.Method {

	case "GET":
		var kittens []Kitten

		err = c.Find(bson.M{}).All(&kittens)
		if err != nil {
			handleError(w, err)
		}

		w.Header().Set("Content-type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(kittens)
		if err != nil {
			handleError(w, err)
		}

	case "POST":
		var kitten Kitten

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handleError(w, err)
		}
		err = r.Body.Close()
		if err != nil {
			handleError(w, err)
		}
		err = json.Unmarshal(body, &kitten)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		err = c.Insert(kitten)
		if err != nil {
			handleError(w, err)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Kitten '%s' created", kitten.Name)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// indexHandler returns a simple message
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I am awesome!")
}

// handleError handles fatal errors
func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Fatalln(err.Error())
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	if vcapServices := os.Getenv("VCAP_SERVICES"); len(vcapServices) == 0 {
		dbString = "localhost"
	} else {
		appEnv, err := cfenv.Current()
		if err != nil {
			log.Fatal(err)
		}
		dbService, err := appEnv.Services.WithName("my-mongodb")
		if err != nil {
			log.Fatal(err)
		}
		uri, ok := dbService.Credentials["uri"].(string)
		if !ok {
			log.Fatal("No valid databse URI found")
		}
		dbString = uri
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/kittens", kittenHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
