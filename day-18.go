package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Coord struct {
	x int
	y int
}

type World struct {
	state  map[Coord]rune
	width  int
	height int
}

func (w World) render() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.height; x++ {
			fmt.Printf(string(w.state[Coord{x, y}]))
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (w World) score() int {
	numWood := 0
	numLumber := 0
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.height; x++ {
			val := w.state[Coord{x, y}]
			if val == '|' {
				numWood++
			} else if val == '#' {
				numLumber++
			}
		}
	}
	return numWood * numLumber
}

func (w World) iterate() World {
	newState := make(map[Coord]rune)
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.height; x++ {
			coords := Coord{x, y}
			current := w.state[coords]
			newState[coords] = current
			if current == '.' {
				numTrees := 0
				for i := x - 1; i <= x+1; i++ {
					for j := y - 1; j <= y+1; j++ {
						if i == x && j == y {
							continue
						}
						if w.state[Coord{i, j}] == '|' {
							numTrees++
						}
					}
				}
				if numTrees >= 3 {
					newState[coords] = '|'
				}
			} else if current == '|' {
				numLumber := 0
				for i := x - 1; i <= x+1; i++ {
					for j := y - 1; j <= y+1; j++ {
						if i == x && j == y {
							continue
						}
						if w.state[Coord{i, j}] == '#' {
							numLumber++
						}
					}
				}
				if numLumber >= 3 {
					newState[coords] = '#'
				}
			} else {
				numLumber := 0
				numTrees := 0
				for i := x - 1; i <= x+1; i++ {
					for j := y - 1; j <= y+1; j++ {
						if i == x && j == y {
							continue
						}
						contents := w.state[Coord{i, j}]
						if contents == '#' {
							numLumber++
						} else if contents == '|' {
							numTrees++
						}
					}
				}
				if numLumber < 1 || numTrees < 1 {
					newState[coords] = '.'
				}
			}
		}
	}
	return World{newState, w.width, w.height}
}

func (w World) hash() string {
	result := ""
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.height; x++ {
			result += string(w.state[Coord{x, y}])
		}
	}
	return result
}

func main() {
	buf, _ := ioutil.ReadFile("day-18.txt")
	input := string(buf)
	lines := strings.Split(input, "\n")
	state := make(map[Coord]rune)
	for y, line := range lines {
		for x, char := range line {
			state[Coord{x, y}] = char
		}
	}
	var world = World{state, len(lines[0]), len(lines)}
	iterations := 10000
	worldMemo := make(map[string]int)
	for i := 0; i < iterations; i++ {
		if i == 10 {
			fmt.Println("Part One:", world.score())
		}
		world = world.iterate()
		existingHash := worldMemo[world.hash()]
		if existingHash > 0 {
			period := i - existingHash
			if i%period == (1000000000-1)%period {
				fmt.Println("Part Two:", world.score())
				break
			}
		}
		worldMemo[world.hash()] = i
	}
}
