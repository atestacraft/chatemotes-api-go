package services

import (
	resoucepack "chatemotes/internal/resourcepack"

	db "github.com/sonyarouje/simdb"
)

type Services struct {
	ResoucePack *resoucepack.ResourcePack
	Database    *db.Driver
}
