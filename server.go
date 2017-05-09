package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Word struct {
	Name string
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		mongoURL := os.Getenv("MONGODB_URI")
		if mongoURL == "" {
			log.Fatal("$MONGODB_URI must be set")
		}

		session, err := mgo.Dial(mongoURL)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()

		c := session.DB("heroku_n1ctsmn2").C("word")
		err = c.Insert(&Word{"negligence",}, &Word{"vulnerabilities",}, &Word{"drawback",}, &Word{"misuse",}, &Word{"damage",})
		if err != nil {
			log.Fatal(err)
		}

		result := Word{}
		err = c.Find(bson.M{"name":"negligence"}).One(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, result.Name)
	})

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)
	n.Run(":"+port)
}
