package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Kitten is a cute kitty
type Kitten struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name string        `json:"name"`
}

// Kittens is an array of kittens
type Kittens []Kitten

// Server is a server with an mgo session
type Server struct {
	db *mgo.Session
}

// KittenHandler allows to manage kittens
func (s *Server) KittenHandler(w http.ResponseWriter, r *http.Request) {
	sess := s.db.Clone()
	defer sess.Close()
	c := sess.DB("").C("kittens")

	switch r.Method {

	case "GET":
		var kittens Kittens

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

		err := json.NewDecoder(r.Body).Decode(&kitten)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = c.Insert(kitten)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, "Kitten '%s' created", kitten.Name)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// IndexHandler returns a simple message
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I am awesome!")
}

func main() {
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
			log.Fatal("no valid databse URI found")
		}
		dbString = uri
	}
	session, err := mgo.Dial(dbString)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	s := Server{db: session}

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/kittens", s.KittenHandler)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
