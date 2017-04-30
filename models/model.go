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

type Animal struct {
	ID          int    `json:"id" sql:"id,pk"`
	Name        string `json:"name" sql:"animalname"`
	Description string `json:"description" sql:"animaldescription"`
	ImageURL    string `json:"image_url" sql:"animalimageurl"`
	Warning     int    `json:"warning" sql:"animaldangerwarning"`
	CategoryID  int    `json:"category_id" sql:"animalcategoryid"`
}

type Sighting struct {
	ID                 int       `json:"id" sql:"sightingid,pk"`
	CreatedAt          time.Time `json:"created_at" sql:"sightingtime"`
	Latitude           float64   `json:"latitude" sql:"latitude"`
	Longitude          float64   `json:"longitude" sql:"longitude"`
	ImageURL           string    `json:"image_url" sql:"sightingimageurl"`
	OriginalSightingID int       `json:"original_sighting_id" sql:"originalsightingid"`
	Gone               bool      `json:"gone" sql:"isanimalgone"`
	ParkID             int       `json:"-" sql:"parkid"`
	UserID             int       `json:"-" sql:"userid"`
	AnimalID           int       `json:"-" sql:"animal_id"`
	Animal             Animal    `json:"animal"`
}

func (m *Model) GetSightings() ([]Sighting, error) {
	var sightings []Sighting
	if err := m.Model(&sightings).OrderExpr("sightingtime asc").Select(); err != nil {
		return nil, err
	}

	return sightings, nil
}

func (m *Model) GetSighting(id int) (*Sighting, error) {
	sighting := Sighting{ID: id}
	if err := m.Model(&sighting).Column("sighting.*", "Animal").Where("sightingid = ?", id).Select(); err != nil {
		return nil, err
	}

	return &sighting, nil
}

func (m *Model) CreateSighting(s Sighting) error {
	return m.Insert(&s)
}

func (m *Model) GetAnimalIDFromName(name string) (int, error) {
	animal := Animal{Name: name}
	if err := m.Model(&animal).Where("animalname = ?", name).Select(); err != nil {
		return 0, err
	}

	return animal.ID, nil
}

func NewModel() *Model {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db := pg.Connect(options)

	return &Model{db}
}
