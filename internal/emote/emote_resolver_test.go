package emote_resolver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	emote_resolver "chatemotes/internal/emote"
)

func TestResolveUrl_ok(t *testing.T) {
	for url, want := range map[string]string{
		"https://7tv.app/emotes/60ae958e229664e8667aea38":                  "https://cdn.7tv.app/emote/60ae958e229664e8667aea38/2x.webp",
		"https://cdn.7tv.app/emote/60ae958e229664e8667aea38/4x.webp":       "https://cdn.7tv.app/emote/60ae958e229664e8667aea38/2x.webp",
		"https://betterttv.com/emotes/5f1b0186cf6d2144653d2970":            "https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/2x",
		"https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/3x.webp": "https://cdn.betterttv.net/emote/5f1b0186cf6d2144653d2970/2x",
		"https://www.frankerfacez.com/emoticon/128054-OMEGALUL":            "https://cdn.frankerfacez.com/emoticon/128054/2",
		"https://cdn.frankerfacez.com/emoticon/128054/3":                   "https://cdn.frankerfacez.com/emoticon/128054/2",
	} {
		t.Run(url, func(t *testing.T) {
			got, ok := emote_resolver.EmoteResolver.ResolveUrl(url)
			assert.True(t, ok)
			assert.Equal(t, want, got)
		})
	}
}
