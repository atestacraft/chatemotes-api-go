package resourcepack

import (
	"archive/zip"
	emote_resolver "chatemotes/internal/emote"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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
		if err := os.Mkdir("pack", 0755); err != nil {
			log.Fatal(err.Error())
		}
	}

	resoucePackFile, err := os.OpenFile(path.Join("pack", "resourcepack.zip"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err.Error())
	}

	resoucePack := &ResourcePack{
		database:         database,
		ResourcePackFile: resoucePackFile,
	}

	fileStat, err := resoucePackFile.Stat()
	if err != nil {
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
	}

	bytes, err := json.Marshal(&ResourcePackMeta{
		Description: "Chat Emotes",
		PackFormat:  9,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	writer.Close()
}

func (r *ResourcePack) GetHash() string {
	hash := sha256.New()
	if _, err := io.Copy(hash, r.ResourcePackFile); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func FetchEmoteImage(url string) (string, error) {
	log.Println("fetching image", url)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes), nil
}

func (r *ResourcePack) AddEmote(url string, name string) (*Emote, error) {
	writer := r.createWriter()
	defer writer.Close()

	file, err := writer.Create(fmt.Sprintf("assets/minecraft/textures/font/%s.png", name))
	if err != nil {
		return nil, err
	}

	emoteUrl, ok := emote_resolver.EmoteResolver.ResolveUrl(url)
	if !ok {
		return nil, errors.New("no match found")
	}

	emoteBase64, err := FetchEmoteImage(emoteUrl)
	if err != nil {
		return nil, err
	}

	emote := &Emote{
		Name:   name,
		Type:   "bitmap",
		File:   fmt.Sprintf("minecraft:font/%s.png", name),
		Height: 10,
		Ascent: 7,
		Chars:  []string{"🤙"},
		Image:  emoteBase64,
	}

	_, err = file.Write([]byte(emoteBase64))
	if err != nil {
		return emote, err
	}

	if err := r.database.Insert(emote); err != nil {
		return emote, err
	}

	return emote, nil
}

func (r *ResourcePack) GetEmotes() ([]Emote, error) {
	var fetchedEmotes []Emote
	err := r.database.
		Open(Emote{}).
		Get().
		AsEntity(&fetchedEmotes)
	if err != nil && err.Error() != "record not found" {
		return nil, xerr.NewW(err)
	}

	return fetchedEmotes, nil
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
