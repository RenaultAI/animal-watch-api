package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"personal/animal-watch-api/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "8080"
)

var m = models.NewModel()

func getSightings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sightings, err := m.GetSightings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	j, err := json.Marshal(sightings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(j))
}

func getSighting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	sightings, err := m.GetSighting(iid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	j, err := json.Marshal(sightings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(j))
}

func main() {
	router := httprouter.New()
	router.GET("/sightings", getSightings)
	router.GET("/sightings/:id", getSighting)

	host := defaultHost
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}
	port := defaultPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	bind := fmt.Sprintf("%s:%s", host, port)

	log.Printf("listening on %s...\n", bind)
	log.Fatal(http.ListenAndServe(bind, router))
}
