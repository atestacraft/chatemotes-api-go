package emote_resolver

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type EmoteResolver struct {
	regex  *regexp.Regexp
	index  int
	imgFmt string
}

type Resolver struct {
	emoteResolver []EmoteResolver
}

func (r EmoteResolver) resolve(url string) (string, bool) {
	matches := r.regex.FindStringSubmatch(url)

	if len(matches) == 0 {
		return "", false
	}

	return fmt.Sprintf(r.imgFmt, matches[r.index]), true
}

func New() *Resolver {
	return &Resolver{
		emoteResolver: []EmoteResolver{
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
	for _, resolver := range r.emoteResolver {
		emoteUrl, ok := resolver.resolve(url)
		if ok {
			return emoteUrl, true
		}
	}

	return "", false
}

func (r *Resolver) FetchEmoteImage(url string) (string, error) {
	log.Println("fetching image", url)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes), nil
}
