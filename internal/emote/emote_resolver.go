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
	matchUrl func(url string) (string, bool)
}

type Resolver struct {
	emoteResolver []EmoteResolver
}

func New() *Resolver {
	return &Resolver{
		emoteResolver: []EmoteResolver{
			// 7tv
			{
				matchUrl: func(url string) (string, bool) {
					r := regexp.MustCompile(`^https://7tv.app/emotes/(\w+)$`)
					matches := getUrlMatches(url, r, 1)
					return fmt.Sprintf("https://cdn.7tv.app/emote/%s/2x.webp", matches), matches != ""
				},
			},
			{
				matchUrl: func(url string) (string, bool) {
					r := regexp.MustCompile(`^https://cdn.7tv.app/emote/(\w+)`)
					matches := getUrlMatches(url, r, 0)
					return fmt.Sprintf("%s/2x.webp", matches), matches != ""
				},
			},
			// bttv
			{
				matchUrl: func(url string) (string, bool) {
					r := regexp.MustCompile(`^https:\/\/betterttv.com\/emotes\/(\w+)`)
					matches := getUrlMatches(url, r, 1)
					return fmt.Sprintf("https://cdn.betterttv.net/emote/%s/2x", matches), matches != ""
				},
			},
			{
				matchUrl: func(url string) (string, bool) {
					r := regexp.MustCompile(`^https:\/\/cdn.betterttv.net\/emote\/(\w+)`)
					matches := getUrlMatches(url, r, 0)
					return fmt.Sprintf("%s/2x", matches), matches != ""
				},
			},
			// ffz
			{
				matchUrl: func(url string) (string, bool) {
					r := regexp.MustCompile(`^https:\/\/www.frankerfacez.com\/emoticon\/(\d+)`)
					matches := getUrlMatches(url, r, 1)
					return fmt.Sprintf("https://cdn.frankerfacez.com/emoticon/%s/2", matches), matches != ""
				},
			},
			{
				matchUrl: func(url string) (string, bool) {
					r := regexp.MustCompile(`^https:\/\/cdn.frankerfacez.com\/emoticon\/(\w+)`)
					matches := getUrlMatches(url, r, 0)
					return fmt.Sprintf("%s/2", matches), matches != ""
				},
			},
		},
	}
}

func getUrlMatches(url string, r *regexp.Regexp, index int) string {
	matches := r.FindStringSubmatch(url)

	if len(matches) == 0 {
		return ""
	}

	return matches[index]
}

func (r *Resolver) ResolveUrl(url string) (string, bool) {
	for _, resolver := range r.emoteResolver {
		emoteUrl, ok := resolver.matchUrl(url)
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
