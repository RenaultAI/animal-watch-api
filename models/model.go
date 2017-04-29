package models

import (
	"log"
	"os"
	"time"

	"github.com/go-pg/pg"
)

type Model struct {
	*pg.DB
}

// type User struct {
// Id     int64
// Name   string
// Emails []string
// }

type Sighting struct {
	ID                 int       `json:"id" sql:"sightingid,pk"`
	CreatedAt          time.Time `json:"created_at" sql:"sightingtime"`
	Latitude           float64   `json:"latitude" sql:"latitude"`
	Longitude          float64   `json:"longitude" sql:"longitude"`
	ImageURL           string    `json:"image_url" sql:"sightingimageurl"`
	OriginalSightingID int       `json:"original_sighting_id" sql:"originalsightingid"`
	Gone               bool      `json:"gone" sql:"isanimalgone"`
	ParkID             int       `json:"park_id" sql:"parkid"`
	AnimalID           int       `json:"animal_id" sql:"animalid"`
	UserID             int       `json:"user_id" sql:"userid"`
}

func (m *Model) GetSightings() ([]Sighting, error) {
	var sightings []Sighting
	if err := m.Model(&sightings).Select(); err != nil {
		return nil, err
	}

	return sightings, nil
}

func (m *Model) GetSighting(id int) (*Sighting, error) {
	sighting := Sighting{ID: id}
	if err := m.Select(&sighting); err != nil {
		return nil, err
	}

	return &sighting, nil
}

func NewModel() *Model {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db := pg.Connect(options)

	return &Model{db}
}
