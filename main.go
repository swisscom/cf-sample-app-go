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

// Controller is a server with an mgo session
type Controller struct {
	session *mgo.Session
}

// kittenHandler allows to manage kittens
func (controller *Controller) kittenHandler(w http.ResponseWriter, r *http.Request) {
	session := controller.session.Clone()
	c := session.DB("").C("kittens")

	switch r.Method {

	case "GET":
		var kittens []Kitten

		err := c.Find(bson.M{}).All(&kittens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(kittens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "POST":
		var kitten Kitten

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &kitten)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		err = c.Insert(kitten)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	var dbString string
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
	session, err := mgo.Dial(dbString)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := Controller{session: session}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/kittens", c.kittenHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
