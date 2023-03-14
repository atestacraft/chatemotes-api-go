package emote_resolver_test

import (
	emote_resolver "chatemotes/internal/emote"
	"testing"

	"github.com/stretchr/testify/assert"
)

var emoteResolver = emote_resolver.New()

func TestSevenTv(t *testing.T) {
	assert := assert.New(t)

	url1, ok := emoteResolver.ResolveUrl("https://7tv.app/emotes/60ae958e229664e8667aea38")
	assert.True(ok)
	assert.Equal(url1, "https://cdn.7tv.app/emote/60ae958e229664e8667aea38/2x.webp")

	url2, ok := emoteResolver.ResolveUrl("https://cdn.7tv.app/emote/60ae958e229664e8667aea38/4x.webp")
	assert.True(ok)
	assert.Equal(url2, "https://cdn.7tv.app/emote/60ae958e229664e8667aea38/2x.webp")
}

func TestBetterTv(t *testing.T) {
	assert := assert.New(t)

	url1, err := emoteResolver.ResolveUrl("https://betterttv.com/emotes/5f1b0186cf6d2144653d2970")
	assert.True(err)
	assert.Equal(url1, "https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/2x")

	url2, err := emoteResolver.ResolveUrl("https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/3x.webp")
	assert.True(err)
	assert.Equal(url2, "https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/2x")
}

func TestFrankerFaceZ(t *testing.T) {
	assert := assert.New(t)

	url1, err := emoteResolver.ResolveUrl("https://www.frankerfacez.com/emoticon/128054-OMEGALUL")
	assert.True(err)
	assert.Equal(url1, "https://cdn.frankerfacez.com/emoticon/128054/2")

	url2, err := emoteResolver.ResolveUrl("https://cdn.frankerfacez.com/emoticon/128054/3")
	assert.True(err)
	assert.Equal(url2, "https://cdn.frankerfacez.com/emoticon/128054/2")
}
