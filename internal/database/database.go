package database

import (
	"github.com/rprtr258/xerr"
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

func (r DB) GetEmotes() ([]Emote, error) {
	return read(r, func(Emote) bool { return true })
}

func (r DB) GetEmoteByName(name string) (*Emote, error) {
	res, err := read(r, func(emote Emote) bool {
		return emote.Name == name
	})
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
