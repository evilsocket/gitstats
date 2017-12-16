/*
 * gitstats - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package main

import "fmt"

// https://misc.flogisoft.com/bash/tip_colors_and_formatting
const (
	BOLD = "\033[1m"
	DIM  = "\033[2m"

	FG_BLACK  = "\033[30m"
	FG_WHITE  = "\033[97m"
	FG_RED    = "\033[31m"
	FG_YELLOW = "\033[33m"
	FG_GREEN  = "\033[32m"

	BG_DGRAY  = "\033[100m"
	BG_RED    = "\033[41m"
	BG_GREEN  = "\033[42m"
	BG_YELLOW = "\033[43m"

	RESET = "\033[0m"
)

func Wrap(s, effect string) string {
	return effect + s + RESET
}

func Dim(s string) string {
	return Wrap(s, DIM)
}

func Bold(s string) string {
	return Wrap(s, BOLD)
}

func Error(s string) string {
	return Wrap(Wrap(s, BG_RED), FG_WHITE)
}

func Bar(val, max, max_width int) {
	perc := float64(val) / float64(max)
	// https://stackoverflow.com/questions/39544571/golang-round-to-nearest-0-05
	width := int((float64(max_width) * perc) + 0.5)
	pad := max_width - width

	third := int((float64(max_width) * 0.33333333333))
	colors := make([]string, max_width)

	for i := 0; i < third; i++ {
		colors[i] = FG_GREEN
	}
	for i := third; i < third*2; i++ {
		colors[i] = FG_YELLOW
	}
	for i := third * 2; i < max_width; i++ {
		colors[i] = FG_RED
	}

	for i := 0; i < width; i++ {
		fmt.Printf(colors[i])
		fmt.Printf("█")
		fmt.Printf(RESET)
	}

	for i := 0; i < pad; i++ {
		fmt.Printf(" ")
		// fmt.Printf("░")
	}
}
