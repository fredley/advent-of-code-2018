package main

import (
	"fmt"
	"strings"
)

func canReact(el1 rune, el2 rune) bool {
	return (el1 != el2) && strings.ToLower(string(el1)) == strings.ToLower(string(el2))
}

func react(polymer string) string {
	var newPolymer string
	for _, char := range polymer {
		if len(newPolymer) == 0 {
			newPolymer = string(char)
			continue
		}
		if canReact(char, rune(newPolymer[len(newPolymer)-1])) {
			newPolymer = newPolymer[:len(newPolymer)-1]
		} else {
			newPolymer = newPolymer + string(char)
		}
	}
	return newPolymer
}

func remove(polymer string, removeChar rune) string {
	var newPolymer string
	var removeStr = strings.ToLower(string(removeChar))
	for _, char := range polymer {
		if strings.ToLower(string(char)) != removeStr {
			newPolymer += string(char)
		}
	}
	return newPolymer
}

func main() {
	const input = ``

	fmt.Println("Part One:", len(react(input)))
	var minLength = len(input)
	for i := 97; i <= 122; i++ { // 97-122 are codes for a-z
		var length = len(react(remove(input, rune(i))))
		if length < minLength {
			minLength = length
		}
	}
	fmt.Println("Part 2:", minLength)
}
