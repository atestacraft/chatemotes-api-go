package database

import (
	"fmt"

	"github.com/rprtr258/simpdb"
	"github.com/rprtr258/xerr"
)

type Emote struct {
	Name   string   `json:"name"`  // ID
	Image  string   `json:"image"` // Base64
	Type   string   `json:"type"`
	File   string   `json:"file"`
	Chars  []string `json:"chars"`
	Height int      `json:"height"`
	Ascent int      `json:"ascent"`
}

func (c Emote) ID() string {
	return c.Name
}

type DB struct {
	table *simpdb.Table[Emote]
}

func New() DB {
	db := simpdb.New("db")
	table, err := simpdb.GetTable[Emote](db, "emotes", simpdb.TableConfig{
		Indent: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer table.Flush()

	return DB{
		table: table,
	}
}

func (r DB) GetLastEmoteChar() string {
	emotes := r.table.Sort(func(e1, e2 Emote) bool {
		return e1.Chars[0] < e2.Chars[0]
	}).Max()

	if emotes.Valid {
		return emotes.Value.Chars[0]
	} else {
		return string(rune(0x1f45f))
	}
}

func (r DB) GetEmotes() []Emote {
	return r.table.List().All()
}

func (r DB) GetEmoteByName(name string) (Emote, error) {
	emote := r.table.Get(name)
	if emote.Valid {
		return emote.Value, nil
	} else {
		return emote.Value, xerr.NewM("emote not found")
	}
}

func (r DB) RemoveEmoteByName(name string) error {
	defer r.table.Flush()
	ok := r.table.DeleteByID(name)
	if ok {
		return nil
	} else {
		return xerr.NewM("emote not found")
	}
}

func (r DB) Insert(emote Emote) {
	defer r.table.Flush()
	r.table.Upsert(emote)
}

func (r DB) UpdateEmote(name, newName string) {
	defer r.table.Flush()
	r.table.Where(func(id string, emote Emote) bool {
		return emote.Name == name
	}).Update(func(emote Emote) Emote {
		emote.Name = newName
		return emote
	})
}
