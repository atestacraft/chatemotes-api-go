package emote_resolver

import (
	"fmt"
	"regexp"
)

type emoteResolver struct {
	regex  *regexp.Regexp
	index  int
	imgFmt string
}

type Resolver struct {
	resolvers []emoteResolver
}

func (r emoteResolver) resolve(url string) (string, bool) {
	matches := r.regex.FindStringSubmatch(url)

	if len(matches) == 0 {
		return "", false
	}

	return fmt.Sprintf(r.imgFmt, matches[r.index]), true
}

var EmoteResolver = newEmoteResolver()

func newEmoteResolver() *Resolver {
	return &Resolver{
		resolvers: []emoteResolver{
			// 7tv
			{
				regex:  regexp.MustCompile(`^https://7tv.app/emotes/(\w+)$`),
				index:  1,
				imgFmt: "https://cdn.7tv.app/emote/%s/2x.webp",
			},
			{
				regex:  regexp.MustCompile(`^https://cdn.7tv.app/emote/(\w+)`),
				index:  0,
				imgFmt: "%s/2x.webp",
			},
			// bttv
			{
				regex:  regexp.MustCompile(`^https:\/\/betterttv.com\/emotes\/(\w+)`),
				index:  1,
				imgFmt: "https://cdn.betterttv.net/emote/%s/2x",
			},
			{
				regex:  regexp.MustCompile(`^https:\/\/cdn.betterttv.net\/emote\/(\w+)`),
				index:  0,
				imgFmt: "%s/2x",
			},
			// ffz
			{
				regex:  regexp.MustCompile(`^https:\/\/www.frankerfacez.com\/emoticon\/(\d+)`),
				index:  1,
				imgFmt: "https://cdn.frankerfacez.com/emoticon/%s/2",
			},
			{
				regex:  regexp.MustCompile(`^https:\/\/cdn.frankerfacez.com\/emoticon\/(\w+)`),
				index:  0,
				imgFmt: "%s/2",
			},
		},
	}
}

func (r *Resolver) ResolveUrl(url string) (string, bool) {
	for _, resolver := range r.resolvers {
		emoteUrl, ok := resolver.resolve(url)
		if ok {
			return emoteUrl, true
		}
	}

	return "", false
}
