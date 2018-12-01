package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	const input = ``

	var freq int64 = 0
	var first_run = true
	var splits = strings.Split(input, "\n")
	var seen = make(map[int64]bool)
	var has_found = false
	for {
		for i, s := range splits {
			var a int64 = 0
			a, _ = strconv.ParseInt(strings.Trim(s, " \t"), 10, 64)
			freq += a
			if seen[freq] && !has_found && i != 0 {
				fmt.Println("Seen (part 2):", freq)
				has_found = true
				if !first_run {
					break
				}
			}
			seen[freq] = true
		}
		if first_run {
			fmt.Println("Frequency (part 1):", freq)
			first_run = false
		}
		if has_found {
			break
		}
	}
}
