package pack

import (
	"archive/zip"
	"chatemotes/internal/database"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rprtr258/xerr"
)

type Pack struct {
	db       database.DB
	filename string
	isFresh  bool
}

func New(resourcePackFile string, db database.DB) *Pack {
	return &Pack{
		db:       db,
		filename: resourcePackFile,
		isFresh:  false,
	}
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

func writeMetadata(w *zip.Writer) error {
	file, err := w.Create("pack.mcmeta")
	if err != nil {
		return err
	}

	_, err = file.Write(metadataBytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *Pack) regenerateFile() error {
	if _, err := os.Stat("pack"); os.IsNotExist(err) {
		if err := os.Mkdir("pack", 0755); err != nil {
			return xerr.NewW(err)
		}
	}

	file, err := os.OpenFile(r.filename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return xerr.NewWM(err, "can't open pack file")
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	if err := writeMetadata(w); err != nil {
		return err
	}

	emotes, err := r.db.GetEmotes()
	if err != nil {
		return xerr.NewWM(err, "can't get emotes list")
	}

	for _, emote := range emotes {
		file, err := w.Create(fmt.Sprintf(
			"assets/minecraft/textures/font/%s.png",
			emote.Name,
		))
		if err != nil {
			return xerr.NewW(err)
		}

		imageBytes, err := base64.RawStdEncoding.DecodeString(emote.Image)
		if err != nil {
			return xerr.NewWM(err, "can't decode image in db",
				xerr.Field("name", emote.Name))
		}

		if _, err = file.Write(imageBytes); err != nil {
			return xerr.NewW(err)
		}
	}

	r.isFresh = true
	return nil
}

func (r *Pack) update() error {
	if r.isFresh {
		return nil
	}

	return r.regenerateFile()
}

func (r *Pack) Invalidate() {
	r.isFresh = false
}

func (r *Pack) Hash() (string, error) {
	if err := r.update(); err != nil {
		return "", err
	}

	file, err := os.OpenFile(r.filename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return "", xerr.NewWM(err, "can't open pack file")
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
