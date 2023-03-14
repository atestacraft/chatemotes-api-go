package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rprtr258/xerr"
)

type entity interface {
	ID() string
}

func entityName[E entity]() string {
	var e E
	return reflect.TypeOf(e).Name()
}

type Emote struct {
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Type   string   `json:"type"`
	File   string   `json:"file"`
	Height int      `json:"height"`
	Ascent int      `json:"ascent"`
	Chars  []string `json:"chars"`
}

func (c Emote) ID() string {
	return c.File
}

type DB struct {
	dir string
}

func New() DB {
	return DB{
		dir: "db",
	}
}

func read[E entity](r DB, filter func(E) bool) ([]E, error) {
	bytes, err := os.ReadFile(filepath.Join(r.dir, entityName[E]()))
	if err != nil {
		return nil, xerr.NewWM(err, "can't open table file",
			xerr.Field("entity", entityName[E]()))
	}

	var all []E
	if err := json.Unmarshal(bytes, &all); err != nil {
		return nil, xerr.NewW(err)
	}

	res := make([]E, 0, len(all))
	for _, entity := range all {
		if filter(entity) {
			res = append(res, entity)
		}
	}

	return res, nil
}

func write[E entity](r DB, entities []E) error {
	bytes, err := json.Marshal(entities)
	if err != nil {
		return xerr.NewW(err)
	}

	if err := os.WriteFile(filepath.Join(r.dir, entityName[E]()), bytes, 0644); err != nil {
		return xerr.NewW(err)
	}

	return nil
}

func (r DB) GetEmotes() ([]Emote, error) {
	return read(r, func(Emote) bool { return true })
}

func (r DB) GetEmoteByName(name string) (*Emote, error) {
	res, err := read(r, func(Emote) bool { return true })
	if err != nil {
		return nil, xerr.NewW(err)
	}

	switch len(res) {
	case 0:
		return nil, nil
	case 1:
		return &res[0], nil
	default:
		return nil, xerr.NewM("too many records",
			xerr.Field("name", name),
			xerr.Field("len(result)", len(res)),
		)
	}
}

func (r DB) RemoveEmoteByName(name string) error {
	all, err := read(r, func(Emote) bool { return true })
	if err != nil {
		return xerr.NewW(err)
	}

	res := make([]Emote, 0, len(all))
	for _, emote := range all {
		if emote.Name != name {
			res = append(res, emote)
		}
	}

	return write(r, res)
}

func (r DB) Insert(emote Emote) error {
	all, err := read(r, func(Emote) bool { return true })
	if err != nil {
		return xerr.NewW(err)
	}

	return write(r, append(all, emote))
}

func (r DB) UpdateEmote(name string) (Emote, error) {
	// TODO: чего апдейт то нахуй
	// поменять Name эмоута с Name==name на name?
	return Emote{}, xerr.NewM("not implemented")
}
