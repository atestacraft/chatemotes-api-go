package logic

import (
	"archive/zip"
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
	"path/filepath"

	"chatemotes/internal/database"
	"chatemotes/internal/emote_resolver"
)

type Logic struct {
	database         database.DB
	ResourcePackFile *os.File
}

func New(database database.DB) Logic {
	if _, err := os.Stat("pack"); os.IsNotExist(err) {
		if err := os.Mkdir("pack", 0755); err != nil {
			log.Fatal(err.Error())
		}
	}

	resoucePackFile, err := os.OpenFile(filepath.Join("pack", "resourcepack.zip"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err.Error())
	}

	resoucePack := Logic{
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

func (r *Logic) createWriter() *zip.Writer {
	return zip.NewWriter(r.ResourcePackFile)
}

var metadataBytes = getMetadataBytes()

func getMetadataBytes() []byte {
	type metadata struct {
		Description string `json:"description"`
		PackFormat  int    `json:"pack_format"`
	}
	bytes, err := json.Marshal(metadata{
		Description: "Chat Emotes",
		PackFormat:  9,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return bytes
}

func (r *Logic) addMetadata() {
	writer := r.createWriter()
	defer writer.Close()

	file, err := writer.Create("pack.mcmeta")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = file.Write(metadataBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (r *Logic) GetHash() string {
	hash := sha256.New()
	if _, err := io.Copy(hash, r.ResourcePackFile); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func downloadImage(url string) ([]byte, error) {
	log.Println("fetching image", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *Logic) AddEmote(url string, name string) (*database.Emote, error) {
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

	imageBytes, err := downloadImage(emoteUrl)
	if err != nil {
		return nil, err
	}

	emoteBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes)

	emote := &database.Emote{
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

func (r *Logic) GetEmotes() ([]database.Emote, error) {
	return r.database.GetEmotes()
}

func (r *Logic) UpdateEmote(name string) (database.Emote, error) {
	return r.database.UpdateEmote(name)
}

func (r *Logic) GetEmoteByName(name string) (*database.Emote, error) {
	return r.database.GetEmoteByName(name)
}

func (r *Logic) RemoveEmoteByName(name string) error {
	return r.database.RemoveEmoteByName(name)
}
