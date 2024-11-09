package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	punctRegexp              = regexp.MustCompile(`^[[:punct:]]|[[:punct:]]$`)
	punctRegexpBegining      = regexp.MustCompile(`^[[:punct:]]`)
	punctRegexpEnd           = regexp.MustCompile(`[[:punct:]]$`)
	punctRegexpMultiBegining = regexp.MustCompile(`^[[:punct:]]{2,}`)
	punctRegexpMultiEnd      = regexp.MustCompile(`[[:punct:]]{2,}$`)
)

type mapEntry struct {
	word  string
	count uint32
}

// Make token lowercase and trim punctuation on edges.
// However, we need to preserve original form in case of multiple
// punctuation characters. E.x.  `hello-------` is considered as a valid
// word, because according to specification `-------` is
// postulated to be a word and we can't trim its edges.  That means we may
// consider entities like `,,,` or `,#!?` as valid words as well. In meanwhile
// `,hello,,,` should be transformed into `hello,,,`.
func adjustWord(w string) string {
	w = strings.ToLower(w)
	// do nothing for `,,hello,,` case
	switch {
	case punctRegexpMultiBegining.MatchString(w) && punctRegexpMultiEnd.MatchString(w):
		return w
	case punctRegexpMultiBegining.MatchString(w):
		return punctRegexpEnd.ReplaceAllString(w, "")
	case punctRegexpMultiEnd.MatchString(w):
		return punctRegexpBegining.ReplaceAllString(w, "")
	default:
		return punctRegexp.ReplaceAllString(strings.ToLower(w), "")
	}
}

// Simple sorting algorithm. First calculate frequencies via hash map, then
// sort (key, value) pairs.
func Top10(str string) []string {
	if len(str) == 0 {
		return []string{}
	}

	counterMap := make(map[string]uint32)
	tokens := strings.Fields(str)

	for _, token := range tokens {
		word := adjustWord(token)
		if len(word) == 0 {
			continue
		}
		counterMap[word]++
	}

	if len(counterMap) == 0 {
		return []string{}
	}

	mapList := make([]mapEntry, 0, len(counterMap))

	for word, count := range counterMap {
		mapList = append(mapList, mapEntry{word, count})
	}

	sort.Slice(mapList, func(i int, j int) bool {
		if mapList[i].count == mapList[j].count {
			return mapList[i].word < mapList[j].word
		}
		return mapList[i].count > mapList[j].count
	})

	result := make([]string, 0, 10)

	for i := 0; i < len(mapList) && i < 10; i++ {
		result = append(result, mapList[i].word)
	}

	return result
}
