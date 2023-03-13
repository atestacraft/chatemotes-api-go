package resourcepack

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type ResourcePackMeta struct {
	Description string `json:"description"`
	PackFormat  int    `json:"pack_format"`
}

type ResourcePack struct {
	Pack ResourcePackMeta `json:"pack"`
}

type Resourcepack struct {
	ResourcePackFile *os.File
	ZipWriter        *zip.Writer
}

func New() *Resourcepack {
	resoucePackFile, err := os.OpenFile("resourcepack.zip", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	zipWriter := zip.NewWriter(resoucePackFile)

	resoucePack := &Resourcepack{
		ResourcePackFile: resoucePackFile,
		ZipWriter:        zipWriter,
	}

	return resoucePack
}

func (r *Resourcepack) GetHash() string {
	hash := sha256.New()
	_, err := io.Copy(hash, r.ResourcePackFile)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (r *Resourcepack) AddEmote() {
}

func (r *Resourcepack) AddMetadata() {
	file, err := r.ZipWriter.Create("pack.mcmeta")
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
}

func (r *Resourcepack) WriteResoucePack() {
	err := r.ZipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Resourcepack) ReadResoucePack() {

}
