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
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"sync"
)

const MAX_WIDTH = 80

type Stats struct {
	lock *sync.Mutex

	Analyzed int
	Skipped  int
	ByUser   map[string]int
	ByHour   []int
	ByDay    []int
	ByMonth  []int
	ByYear   map[int]int
	Tags     map[string]int
}

func NewStats() Stats {
	return Stats{
		lock: &sync.Mutex{},

		Analyzed: 0,
		Skipped:  0,
		ByUser:   make(map[string]int, 0),
		ByHour:   make([]int, 24),
		ByDay:    make([]int, 7),
		ByMonth:  make([]int, 12),
		ByYear:   make(map[int]int, 0),
		Tags:     make(map[string]int, 0),
	}
}

func (s *Stats) Feed(c *object.Commit) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Analyzed++

	who := c.Author
	when := who.When
	what := c.Message

	for _, t := range Tokenize(what) {
		if _, found := s.Tags[t]; found == true {
			s.Tags[t] += 1
		} else {
			s.Tags[t] = 1
		}
	}

	if _, found := s.ByUser[who.Email]; found == true {
		s.ByUser[who.Email] += 1
	} else {
		s.ByUser[who.Email] = 1
	}

	if _, found := s.ByYear[when.Year()]; found == true {
		s.ByYear[when.Year()] += 1
	} else {
		s.ByYear[when.Year()] = 1
	}

	s.ByHour[when.Hour()] += 1
	s.ByDay[when.Weekday()] += 1
	s.ByMonth[when.Month()-1] += 1
}

func (s *Stats) Skip(c *object.Commit) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Skipped++
}

func (s *Stats) pTitle(t string, w int) {
	fmt.Printf("%s\n", Bold(t))
	for i := 0; i < w; i++ {
		fmt.Printf("─")
	}
	fmt.Println()
}

func (s *Stats) pChart(title string, format string, labels []string, values []int) {
	s.pTitle(title, MAX_WIDTH)

	max_txt_len := 0
	texts := make([]string, 0)
	for idx, value := range values {
		label := ""
		text := ""
		if labels != nil {
			label = labels[idx]
		} else {
			label = fmt.Sprintf(format, idx)
		}
		text = fmt.Sprintf(" "+format+" | %s", value, label)

		len := len(text)
		if len > max_txt_len {
			max_txt_len = len
		}
		texts = append(texts, text)
	}

	for idx, value := range values {
		width_left := MAX_WIDTH - max_txt_len
		Bar(value, s.Analyzed, width_left)
		fmt.Printf("%s\n", texts[idx])
	}
}

func (s *Stats) Print() {
	// HACK lulz
	vFormat := "%" + fmt.Sprintf("%d", len(fmt.Sprintf("%d", s.Analyzed))) + "d"

	if len(s.ByUser) == 1 {
		user := ""
		tot := 0
		for u, t := range s.ByUser {
			user = u
			tot = t
		}

		fmt.Printf("Author %s made %d total commits.", Bold(user), tot)
	} else {
		users := Keys(s.ByUser)
		values := Values(s.ByUser, users)
		s.pChart("Per user activity distribution.",
			vFormat,
			users,
			values)
	}

	fmt.Println()

	hours := make([]string, 0)
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d:00", i))
	}

	s.pChart("Daily commits distribution (per hour).",
		vFormat,
		hours,
		s.ByHour)

	fmt.Println()

	s.pChart("Weekly commits distribution (by day).",
		vFormat,
		[]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		s.ByDay)

	fmt.Println()

	s.pChart("Monthly commits distribution.",
		vFormat,
		[]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul",
			"Aug", "Sep", "Oct", "Nov", "Dec"},
		s.ByMonth)

	fmt.Println()

	sortedYears := SortedKeys(s.ByYear)
	sortedValues := make([]int, 0)

	for _, year := range sortedYears {
		sortedValues = append(sortedValues, s.ByYear[year])
	}

	s.pChart("Yearly commits distribution.",
		vFormat,
		Stringize(sortedYears),
		sortedValues)

	fmt.Println()

	top_tags := 25
	tags := Tags(s.Tags, top_tags)
	ntags := len(tags)
	labels := make([]string, ntags)
	values := make([]int, ntags)
	for i, t := range tags {
		labels[i] = t.Value
		values[i] = t.Hits
	}

	s.pChart(fmt.Sprintf("Words distribution (top %d).", top_tags),
		vFormat,
		labels,
		values)
}
