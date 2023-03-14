package logic

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"

	"github.com/rprtr258/xerr"
	"github.com/tidbyt/go-libwebp/webp"

	"chatemotes/internal/database"
	"chatemotes/internal/emote_resolver"
	"chatemotes/internal/pack"
)

type Logic struct {
	db   database.DB
	pack *pack.Pack
}

func New(packFilename string, db database.DB) Logic {
	return Logic{
		db:   db,
		pack: pack.New(packFilename, db),
	}
}

func (r *Logic) GetHash() string {
	hash, err := r.pack.Hash()
	if err != nil {
		log.Fatal(err.Error())
	}
	return hash
}

func downloadImage(url string) ([]byte, error) {
	log.Println("fetching image", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, xerr.NewW(err)
	}
	defer response.Body.Close()

	webpBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, xerr.NewW(err)
	}

	webpDecoder, err := webp.NewAnimationDecoder(webpBytes)
	if err != nil {
		return nil, xerr.NewW(err)
	}
	defer webpDecoder.Close()

	anim, err := webpDecoder.Decode()
	if err != nil {
		return nil, xerr.NewW(err)
	}

	var out bytes.Buffer
	if err := png.Encode(&out, anim.Image[0]); err != nil {
		return nil, xerr.NewW(err)
	}

	return out.Bytes(), nil
}

func (r *Logic) AddEmote(url string, name string) (*database.Emote, error) {
	imageURL, ok := emote_resolver.EmoteResolver.ResolveUrl(url)
	if !ok {
		return nil, errors.New("no match found")
	}

	imageBytes, err := downloadImage(imageURL)
	if err != nil {
		return nil, err
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	emote := &database.Emote{
		Name:   name,
		Type:   "bitmap",
		File:   fmt.Sprintf("minecraft:font/%s.png", name),
		Height: 10,
		Ascent: 7,
		Chars:  []string{"🤙"},
		Image:  "data:image/png;base64," + imageBase64,
	}

	if err := r.db.Insert(emote); err != nil {
		return nil, err
	}

	r.pack.Invalidate()

	return emote, nil
}

func (r *Logic) GetEmotes() ([]database.Emote, error) {
	return r.db.GetEmotes()
}

func (r *Logic) UpdateEmote(name string) (database.Emote, error) {
	return r.db.UpdateEmote(name)
}

func (r *Logic) GetEmoteByName(name string) (*database.Emote, error) {
	return r.db.GetEmoteByName(name)
}

func (r *Logic) RemoveEmoteByName(name string) error {
	return r.db.RemoveEmoteByName(name)
}
