package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"personal/animal-watch-api/models"
	"strconv"
	"time"

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
		return
	}

	// only get latest sightings for distinct sightings
	mapping := make(map[int]int, len(sightings))
	result := make([]models.Sighting, 0, len(sightings))
	index := 0
	for _, s := range sightings {
		id := s.OriginalSightingID
		if id == 0 {
			id = s.ID
		}
		j, exists := mapping[id]
		if exists {
			result[j] = s
		} else {
			result = append(result, s)
			mapping[id] = index
			index++
		}
	}

	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(j))
}

func getSighting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sightings, err := m.GetSighting(iid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(sightings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(j))
}

func createSighting(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sighting models.Sighting
	json.NewDecoder(r.Body).Decode(&sighting)

	if sighting.Latitude == 0 || sighting.Longitude == 0 || sighting.AnimalID == 0 {
		http.Error(w, fmt.Sprintf("Bad parameter: latitude, longitude, animal ID are required."), http.StatusBadRequest)
		return
	}

	sighting.CreatedAt = time.Now()
	sighting.ParkID = 1
	if err := m.CreateSighting(sighting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func optionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	router := httprouter.New()
	router.OPTIONS("/sightings", optionsHandler)
	router.GET("/sightings", getSightings)
	router.GET("/sightings/:id", getSighting)
	router.POST("/sightings", createSighting)

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
