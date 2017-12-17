/*
 * gitstats - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

func URLValid(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	} else {
		return true
	}
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func SortedKeys(m map[int]int) []int {
	keys := make([]int, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func Stringize(ints []int) []string {
	strs := make([]string, 0)
	for _, v := range ints {
		strs = append(strs, fmt.Sprintf("%d", v))
	}
	return strs
}

func Keys(m map[string]int) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values(m map[string]int, keys []string) []int {
	values := make([]int, 0)
	for _, k := range keys {
		values = append(values, m[k])
	}
	return values
}

var tokenizer = regexp.MustCompile("[^\\w]+")

func Tokenize(text string) []string {
	// strip signature
	if strings.Contains(text, "-----END PGP SIGNATURE-----") == true {
		lines := strings.Split(text, "\n")
		cut_at := 0
		for i, line := range lines {
			if strings.Trim(line, "\n\r\t ") == "-----END PGP SIGNATURE-----" {
				cut_at = i
				break
			}
		}

		text = strings.Join(lines[cut_at:], "\n")
	}
	// trim
	text = strings.Trim(text, "\n\r\t ")
	// split words
	tokens := tokenizer.Split(text, -1)
	// make lower, remove empty or short or blacklisted
	tmp := make([]string, 0)
	for _, t := range tokens {
		t = strings.Trim(t, "\n\r\t ")
		t = strings.ToLower(t)
		if t != "" && len(t) >= 3 {
			tmp = append(tmp, t)
		}
	}
	tokens = tmp
	// get unique values
	m := make(map[string]int, 0)
	for i, t := range tokens {
		m[t] = i
	}
	tokens = Keys(m)

	return tokens
}

type Tag struct {
	Value string
	Hits  int
}

func Tags(tags map[string]int, size int) []Tag {
	cloud := make([]Tag, 0)

	for t, h := range tags {
		cloud = append(cloud, Tag{t, h})
	}

	sort.Slice(cloud, func(i, j int) bool {
		return cloud[i].Hits > cloud[j].Hits
	})

	if len(cloud) > size {
		cloud = cloud[0:size]
	}
	return cloud
}
