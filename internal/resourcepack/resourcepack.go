package resourcepack

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/sonyarouje/simdb"
)

type Emote struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	File   string   `json:"file"`
	Height int      `json:"height"`
	Ascent int      `json:"ascent"`
	Chars  []string `json:"chars"`
}

func (c Emote) ID() (jsonField string, value interface{}) {
	value = c.File
	jsonField = "file"
	return
}

type ResourcePackMeta struct {
	Description string `json:"description"`
	PackFormat  int    `json:"pack_format"`
}

type McMeta struct {
	Pack ResourcePackMeta `json:"pack"`
}

type ResourcePack struct {
	database         *simdb.Driver
	ResourcePackFile *os.File
}

func New(database *simdb.Driver) *ResourcePack {
	if _, err := os.Stat("pack"); os.IsNotExist(err) {
		os.Mkdir("pack", 0755)
	}

	resoucePackFile, err := os.OpenFile(path.Join("pack", "resourcepack.zip"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	resoucePack := &ResourcePack{
		database:         database,
		ResourcePackFile: resoucePackFile,
	}

	fileStat, err := resoucePackFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileStat.Size() == 0 {
		resoucePack.addMetadata()
	}

	return resoucePack
}

func (r *ResourcePack) createWriter() *zip.Writer {
	return zip.NewWriter(r.ResourcePackFile)
}

func (r *ResourcePack) addMetadata() {
	writer := r.createWriter()
	file, err := writer.Create("pack.mcmeta")
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := json.Marshal(&ResourcePackMeta{
		Description: "Chat Emotes",
		PackFormat:  9,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

	writer.Close()
}

func (r *ResourcePack) GetHash() string {
	hash := sha256.New()
	_, err := io.Copy(hash, r.ResourcePackFile)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (r *ResourcePack) AddEmote( /*emoteImage []byte,*/ name string) {
	emote := &Emote{
		Name:   name,
		Type:   "bitmap",
		File:   fmt.Sprintf("minecraft:font/%s.png", name),
		Height: 10,
		Ascent: 7,
		Chars:  []string{"🤙"},
	}

	err := r.database.Insert(emote)
	if err != nil {
		log.Fatal(err)
	}

	// writer := r.createWriter()
	// file, err := writer.Create(fmt.Sprintf("assets/minecraft/textures/font/%s.png", name))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = file.Write(emoteImage)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// writer.Close()
}

func (r *ResourcePack) GetEmotes() []Emote {
	var fetchedEmotes []Emote
	err := r.database.
		Open(Emote{}).
		Get().
		AsEntity(&fetchedEmotes)

	if err != nil {
		log.Fatal(err)
	}

	return fetchedEmotes
}

func (r *ResourcePack) GetEmoteByName(name string) (Emote, error) {
	var fetchedEmote Emote
	err := r.database.
		Open(Emote{}).
		Where("name", "=", strings.ToLower(name)).
		First().
		AsEntity(&fetchedEmote)

	return fetchedEmote, err
}
