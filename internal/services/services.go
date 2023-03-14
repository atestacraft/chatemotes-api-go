package services

import (
	emote_resolver "chatemotes/internal/emote"
	resoucepack "chatemotes/internal/resourcepack"

	db "github.com/sonyarouje/simdb"
)

type Services struct {
	EmoteResolver *emote_resolver.Resolver
	ResoucePack   *resoucepack.ResourcePack
	Database      *db.Driver
}
