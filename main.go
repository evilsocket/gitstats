/*
 * gitstats - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package main

import (
	"flag"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"strings"
	"time"
)

const VERSION = "1.0.0"

var (
	stats    = NewStats()
	repo     = (*git.Repository)(nil)
	path     = ""
	sauthors = ""
	authors  = make(map[string]bool, 0)
)

func init() {
	flag.StringVar(&path, "repo", path, "Repository to analyze, can be a URL or a folder.")
	flag.StringVar(&sauthors, "authors", "", "Comma separated list of email addresses to filter commits with.")
}

func fatal(err string) {
	fmt.Printf("%s\n", Error(err))
	os.Exit(1)
}

func main() {
	var err error

	flag.Parse()

	if Exists(path) == false {
		fatal("Invalid path specified")
	}

	emails := strings.Split(sauthors, ",")
	for _, email := range emails {
		email = strings.Trim(email, "\n\r\t ")
		if email != "" {
			authors[email] = true
		}
	}

	started := time.Now()

	fmt.Printf("gitstats v%s is starting the analysis on %s ...\n\n", VERSION, Bold(path))

	repo, err = git.PlainOpen(path)
	if err != nil {
		fatal(err.Error())
	}

	ref, err := repo.Head()
	if err != nil {
		fatal("Could not get HEAD: " + err.Error())
	}

	history, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		fatal("Could not get commits history: " + err.Error())
	}

	err = history.ForEach(func(c *object.Commit) error {
		analyze := true
		if len(authors) > 0 {
			_, analyze = authors[c.Author.Email]
		}

		if analyze {
			stats.Feed(c)
		} else {
			stats.Skip(c)
		}

		return nil
	})

	total_time := time.Since(started)

	fmt.Println()

	stats.Print()

	nauthors := len(stats.ByUser)
	plural := "s"
	if nauthors == 1 {
		plural = ""
	}
	fmt.Printf("\nAnalyzed %d commits (%d skipped) by %d author%s in %s.\n", stats.Analyzed, stats.Skipped, nauthors, plural, total_time)
}
