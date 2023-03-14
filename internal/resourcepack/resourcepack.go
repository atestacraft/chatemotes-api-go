package resourcepack

import (
	"archive/zip"
	emote_resolver "chatemotes/internal/emote"
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
	Image  string   `json:"image"`
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
	emoteResolver    *emote_resolver.Resolver
	ResourcePackFile *os.File
}

func New(database *simdb.Driver, emoteResolver *emote_resolver.Resolver) *ResourcePack {
	if _, err := os.Stat("pack"); os.IsNotExist(err) {
		os.Mkdir("pack", 0755)
	}

	resoucePackFile, err := os.OpenFile(path.Join("pack", "resourcepack.zip"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	resoucePack := &ResourcePack{
		database:         database,
		emoteResolver:    emoteResolver,
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

func (r *ResourcePack) AddEmote(url string, name string) (*Emote, error) {
	emote := &Emote{
		Name:   name,
		Type:   "bitmap",
		File:   fmt.Sprintf("minecraft:font/%s.png", name),
		Height: 10,
		Ascent: 7,
		Chars:  []string{"🤙"},
	}

	writer := r.createWriter()
	file, err := writer.Create(fmt.Sprintf("assets/minecraft/textures/font/%s.png", name))
	if err != nil {
		return emote, err
	}

	emoteUrl, err := r.emoteResolver.ResolveUrl(url)
	if err != nil {
		return emote, err
	}

	emoteBase64, err := r.emoteResolver.FetchEmoteImage(emoteUrl)
	if err != nil {
		return emote, err
	}

	emote.Image = emoteBase64

	_, err = file.Write([]byte(emoteBase64))
	if err != nil {
		return emote, err
	}

	if r.database.Insert(emote); err != nil {
		return emote, err
	}

	writer.Close()

	return emote, nil
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

func (r *ResourcePack) RemoveEmoteByName(name string) error {
	emote := Emote{Name: name}
	err := r.database.
		Open(Emote{}).
		Where("name", "=", strings.ToLower(name)).
		Delete(&emote)

	return err
}

func (r *ResourcePack) UpdateEmote(url string, name string) (Emote, error) {
	emote := Emote{Name: name}
	err := r.database.
		Open(Emote{}).
		Where("name", "=", strings.ToLower(name)).
		Update(&emote)

	return emote, err
}
