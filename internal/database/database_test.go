package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityName(t *testing.T) {
	assert.Equal(t, "Emote", entityName[Emote]())
}
