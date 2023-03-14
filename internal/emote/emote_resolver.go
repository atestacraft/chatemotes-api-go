package emote_resolver

import (
	"errors"
	"fmt"
	"regexp"
)

type EmoteResolver struct {
	matchUrl func(url string) (string, string)
}

type Resolver struct {
	emoteResolver []EmoteResolver
}

func New() *Resolver {
	emoteResolvers := []EmoteResolver{
		// 7tv
		{
			matchUrl: func(url string) (string, string) {
				r := regexp.MustCompile(`^https://7tv.app/emotes/(\w+)$`)
				matches := getUrlMatches(url, r, 1)
				return fmt.Sprintf("https://cdn.7tv.app/emote/%s/2x", matches), matches
			},
		},
		{
			matchUrl: func(url string) (string, string) {
				r := regexp.MustCompile(`^https://cdn.7tv.app/emote/(\w+)`)
				matches := getUrlMatches(url, r, 0)
				return fmt.Sprintf("%s/2x", matches), matches
			},
		},
		// bttv
		{
			matchUrl: func(url string) (string, string) {
				r := regexp.MustCompile(`^https:\/\/betterttv.com\/emotes\/(\w+)`)
				matches := getUrlMatches(url, r, 1)
				return fmt.Sprintf("https://cdn.betterttv.net/emote/%s/2x", matches), matches
			},
		},
		{
			matchUrl: func(url string) (string, string) {
				r := regexp.MustCompile(`^https:\/\/cdn.betterttv.net\/emote\/(\w+)`)
				matches := getUrlMatches(url, r, 0)
				return fmt.Sprintf("%s/2x", matches), matches
			},
		},
		// ffz
		{
			matchUrl: func(url string) (string, string) {
				r := regexp.MustCompile(`^https:\/\/www.frankerfacez.com\/emoticon\/(\d+)`)
				matches := getUrlMatches(url, r, 1)
				return fmt.Sprintf("https://cdn.frankerfacez.com/emoticon/%s/2", matches), matches
			},
		},
		{
			matchUrl: func(url string) (string, string) {
				r := regexp.MustCompile(`^https:\/\/cdn.frankerfacez.com\/emoticon\/(\w+)`)
				matches := getUrlMatches(url, r, 0)
				return fmt.Sprintf("%s/2", matches), matches
			},
		},
	}

	resolver := &Resolver{
		emoteResolver: emoteResolvers,
	}

	return resolver
}

func getUrlMatches(url string, r *regexp.Regexp, index int) string {
	matches := r.FindStringSubmatch(url)

	if len(matches) == 0 {
		return ""
	}

	return matches[index]
}

func (r *Resolver) ResolveUrl(url string) (string, error) {
	for _, resolver := range r.emoteResolver {
		emoteUrl, matches := resolver.matchUrl(url)
		if matches != "" {
			return emoteUrl, nil
		}
	}

	return "", errors.New("no match found")
}
