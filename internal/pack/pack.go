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

const ImageBytesPrefix = "data:image/png;base64,"

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
	type Metadata struct {
		PackFormat  int    `json:"pack_format"`
		Description string `json:"description"`
	}

	type PackMetadata struct {
		Pack Metadata `json:"pack"`
	}

	bytes, err := json.Marshal(&PackMetadata{
		Pack: Metadata{
			PackFormat:  12,
			Description: "Chat Emotes",
		},
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

func (r *Pack) writeFont(w *zip.Writer) error {
	file, err := w.Create("assets/minecraft/font/default.json")
	if err != nil {
		return err
	}

	type Font struct {
		Providers []database.Emote `json:"providers"`
	}

	emotes := r.db.GetEmotes()
	bytes, err := json.Marshal(&Font{Providers: emotes})
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

func (r *Pack) regenerate() error {
	log.Println("regenerating pack")

	if _, err := os.Stat("pack"); err != nil {
		if !os.IsNotExist(err) {
			return xerr.NewW(err)
		}

		if err := os.Mkdir("pack", 0755); err != nil {
			return xerr.NewW(err)
		}
	}

	if err := os.Remove(r.filename); err != nil && !os.IsNotExist(err) {
		return xerr.NewWM(err, "can't remove old pack file")
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

	if err := r.writeFont(w); err != nil {
		return err
	}

	emotes := r.db.GetEmotes()
	for _, emote := range emotes {
		imageContent := emote.Image[len(ImageBytesPrefix):]

		imageBytes, err := base64.StdEncoding.DecodeString(imageContent)
		if err != nil {
			log.Println(xerr.NewWM(err, "can't decode image in db",
				xerr.Field("name", emote.Name)).Error())
			continue
		}

		file, err := w.Create(fmt.Sprintf(
			"assets/minecraft/textures/font/%s.png",
			emote.Name,
		))
		if err != nil {
			return xerr.NewW(err)
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

	return r.regenerate()
}

func (r *Pack) Invalidate() {
	r.isFresh = false
}

func (r *Pack) Filename() string {
	if err := r.update(); err != nil {
		log.Fatal(err.Error())
	}

	return r.filename
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
