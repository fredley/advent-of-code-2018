package main

import (
	"fmt"
	"strings"
)

type Transform struct {
	test   string
	output string
}

func transform(state string, transforms []Transform) string {
	var newState = ".."
	for i := 2; i < len(state)-2; i++ {
		var isTransformed = false
		for _, transform := range transforms {
			if state[i-2:i+3] == transform.test {
				newState += transform.output
				isTransformed = true
				break
			}
		}
		if !isTransformed {
			newState += "."
			// newState += string(state[i])
		}

	}
	return newState + ".."
}

func main() {
	const padding = "......................."
	const initialState = padding + "" + padding
	const input_transforms = ``
	split := strings.Split(input_transforms, "\n")
	transforms := make([]Transform, len(split))
	for i, line := range split {
		line = strings.Trim(line, " ")
		transforms[i] = Transform{line[:5], string(line[9])}
	}
	fmt.Println("Part One:", getSum(20, initialState, transforms, len(padding)))
	var pattern = ""
	var generation = 1
	var state = initialState
	for {
		state = transform(state, transforms) + ".."
		newPattern := getPattern(state)
		if newPattern == pattern {
			baseSum := getSum(generation, initialState, transforms, len(padding))
			fmt.Println("Part Two:", baseSum+((50000000000-generation)*getCount(pattern)))
			break
		} else {
			pattern = newPattern
		}
		generation++
	}
}

func getSum(generations int, state string, transforms []Transform, padding int) int {
	for i := 0; i < generations; i++ {
		state = transform(state, transforms) + ".."
	}
	var sum = 0
	for i := 0; i < len(state); i++ {
		if string(state[i]) == "#" {
			sum += i - padding
		}
	}
	return sum
}

func getPattern(state string) string {
	var startIndex = 0
	var endIndex = 0
	for i := 0; i < len(state); i++ {
		if string(state[i]) == "#" {
			if startIndex == 0 {
				startIndex = i
			} else {
				endIndex = i
			}
		}
	}
	return state[startIndex : endIndex+1]
}

func getCount(state string) int {
	var sum = 0
	for i := 0; i < len(state); i++ {
		if string(state[i]) == "#" {
			sum += 1
		}
	}
	return sum
}
