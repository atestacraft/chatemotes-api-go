package database

import (
	"log"
	"strings"

	"github.com/rprtr258/xerr"
	"github.com/sonyarouje/simdb"
)

type Emote struct {
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Type   string   `json:"type"`
	File   string   `json:"file"`
	Height int      `json:"height"`
	Ascent int      `json:"ascent"`
	Chars  []string `json:"chars"`
}

func (c Emote) ID() (string, any) {
	return "file", c.File
}

type DB struct {
	*simdb.Driver
}

func New() DB {
	driver, err := simdb.New("db")
	if err != nil {
		log.Fatal(err)
	}

	return DB{driver}
}

func (r DB) emotesTable() *simdb.Driver {
	return r.Open(Emote{})
}

func (r DB) GetEmotes() ([]Emote, error) {
	var fetchedEmotes []Emote
	err := r.
		emotesTable().
		Get().
		AsEntity(&fetchedEmotes)
	if err != nil && err.Error() != "record not found" {
		return nil, xerr.NewW(err)
	}

	return fetchedEmotes, nil
}

func (r DB) GetEmoteByName(name string) (*Emote, error) {
	var emote Emote
	err := r.
		emotesTable().
		Where("name", "=", strings.ToLower(name)).
		First().
		AsEntity(&emote)

	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}

		return nil, xerr.NewW(err)
	}

	return &emote, nil
}

func (r DB) RemoveEmoteByName(name string) error {
	err := r.
		emotesTable().
		Where("name", "=", strings.ToLower(name)).
		Delete(&Emote{Name: name})

	return xerr.NewW(err)
}

func (r DB) UpdateEmote(name string) (Emote, error) {
	emote := Emote{Name: name}
	err := r.
		emotesTable().
		Where("name", "=", strings.ToLower(name)).
		Update(&emote)

	return emote, xerr.NewW(err)
}
