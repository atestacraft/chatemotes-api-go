package emote_resolver_test

import (
	emote_resolver "chatemotes/internal/emote"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSevenTv(t *testing.T) {
	url1, ok := emote_resolver.EmoteResolver.ResolveUrl("https://7tv.app/emotes/60ae958e229664e8667aea38")
	assert.True(t, ok)
	assert.Equal(t, url1, "https://cdn.7tv.app/emote/60ae958e229664e8667aea38/2x.webp")

	url2, ok := emote_resolver.EmoteResolver.ResolveUrl("https://cdn.7tv.app/emote/60ae958e229664e8667aea38/4x.webp")
	assert.True(t, ok)
	assert.Equal(t, url2, "https://cdn.7tv.app/emote/60ae958e229664e8667aea38/2x.webp")
}

func TestBetterTv(t *testing.T) {
	url1, ok := emote_resolver.EmoteResolver.ResolveUrl("https://betterttv.com/emotes/5f1b0186cf6d2144653d2970")
	assert.True(t, ok)
	assert.Equal(t, url1, "https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/2x")

	url2, ok := emote_resolver.EmoteResolver.ResolveUrl("https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/3x.webp")
	assert.True(t, ok)
	assert.Equal(t, url2, "https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/2x")
}

func TestFrankerFaceZ(t *testing.T) {
	url1, ok := emote_resolver.EmoteResolver.ResolveUrl("https://www.frankerfacez.com/emoticon/128054-OMEGALUL")
	assert.True(t, ok)
	assert.Equal(t, url1, "https://cdn.frankerfacez.com/emoticon/128054/2")

	url2, ok := emote_resolver.EmoteResolver.ResolveUrl("https://cdn.frankerfacez.com/emoticon/128054/3")
	assert.True(t, ok)
	assert.Equal(t, url2, "https://cdn.frankerfacez.com/emoticon/128054/2")
}
