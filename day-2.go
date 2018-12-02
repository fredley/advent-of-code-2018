package main

import (
	"fmt"
	"strings"
)

func main() {
	const input = ``

	var contains_two int
	var contains_three int
	var seen_words = make([]string, len(strings.Split(input, "\n")))
	var found_part_two = false

	for i, word := range strings.Split(input, "\n") {
		var stripped_word = strings.Trim(word, " \t")
		var counts = make(map[string]int)
		var checksum = 0
		for _, char := range stripped_word {
			counts[string(char)] += 1
			checksum += int(char)
		}
		var found_two = false
		var found_three = false
		for _, count := range counts {
			if count == 2 && !found_two {
				contains_two += 1
				found_two = true
			}
			if count == 3 && !found_three {
				found_three = true
				contains_three += 1
			}
			if found_two && found_three {
				break
			}
		}
		if found_part_two {
			continue
		}
		for j, seen_word := range seen_words {
			if j >= i {
				break
			}
			var diffs = 0
			var diff_char = 0
			for c, char := range stripped_word {
				if int(char) != int(seen_word[c]) {
					diffs += 1
					diff_char = c
				}
			}
			if diffs == 1 {
				found_part_two = true
				fmt.Println("Part two: ", stripped_word[:diff_char]+stripped_word[diff_char+1:])
				break
			}
		}
		seen_words[i] = stripped_word
	}
	fmt.Println("Part one: ", contains_two*contains_three)
}
