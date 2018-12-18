package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type Frontiers [][]int

func (s Frontiers) Len() int {
	return len(s)
}
func (s Frontiers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Frontiers) Less(i, j int) bool {
	l1 := s[i]
	l2 := s[j]
	if l1[1] < l2[1] || (l1[1] == l2[1] && l1[0] < l2[0]) {
		return true
	}
	return false
}

func unique(frontiers [][]int) [][]int {
	sort.Sort(Frontiers(frontiers))
	var result = make([][]int, 0)
	var prev []int
	for _, frontier := range frontiers {
		if len(prev) == 0 || (prev[0] != frontier[0] || prev[1] != frontier[1]) {
			prev = frontier
			result = append(result, frontier)
		}
	}
	return result
}

func worldKey(x int, y int) string {
	return string(x) + "," + string(y)
}

func renderWorld(world map[string]string, minX int, maxX int, minY int, maxY int, frontiers [][]int) {
	for _, frontier := range frontiers {
		world[worldKey(frontier[0], frontier[1])] = "+"
	}
	for j := minY - 1; j <= maxY; j++ {
		for i := minX; i <= maxX+1; i++ {
			char := world[worldKey(i, j)]
			if len(char) == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf(char)
			}

		}
		fmt.Printf("\n")
	}
}

func waterCound(world map[string]string, minX int, maxX int, minY int, maxY int, retainedOnly bool) int {
	var count = 0
	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX+1; i++ {
			char := world[worldKey(i, j)]
			if char == "~" || (char == "|" && !retainedOnly) {
				count++
			}
		}
	}
	return count
}

func main() {
	buf, _ := ioutil.ReadFile("day-17.txt")
	input := string(buf)
	lines := strings.Split(input, "\n")
	world := make(map[string]string)
	var xPart string
	var yPart string
	var minY = 1000
	var maxY = 0
	var minX = 1000
	var maxX = 0
	for _, line := range lines {
		if len(line) == 0 {
			break
		}
		parts := strings.Split(line, ", ")
		if parts[0][0] == 'x' {
			xPart = strings.Split(parts[0], "=")[1]
			yPart = strings.Split(parts[1], "=")[1]
		} else {
			yPart = strings.Split(parts[0], "=")[1]
			xPart = strings.Split(parts[1], "=")[1]
		}
		if strings.Index(xPart, "..") > 1 {
			xParts := strings.Split(xPart, "..")
			x1, _ := strconv.Atoi(xParts[0])
			x2, _ := strconv.Atoi(xParts[1])
			y, _ := strconv.Atoi(yPart)
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
			if x1 < minX {
				minX = x1
			}
			if x2 > maxX {
				maxX = x2
			}
			for i := x1; i <= x2; i++ {
				world[worldKey(i, y)] = "#"
			}
		} else {
			yParts := strings.Split(yPart, "..")
			y1, _ := strconv.Atoi(yParts[0])
			y2, _ := strconv.Atoi(yParts[1])
			x, _ := strconv.Atoi(xPart)
			if y1 < minY {
				minY = y1
			}
			if y2 > maxY {
				maxY = y2
			}
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			for i := y1; i <= y2; i++ {
				world[worldKey(x, i)] = "#"
			}
		}
	}
	var frontiers = [][]int{[]int{500, 0}}
	var clock = 0
	maxClock := 18000
	maxWidClock := maxX - minX
	for {
		clock++
		if clock > maxClock {
			fmt.Println("Stuck")
			fmt.Println(frontiers)
			renderWorld(world, minX, maxX, minY, maxY, frontiers)
			break
		}
		if len(frontiers) == 0 {
			break
		}
		newFrontiers := make([][]int, 0)
		var stillGoing = false
		for _, frontier := range unique(frontiers) {
			x := frontier[0]
			y := frontier[1]
			if y > maxY {
				continue
			}
			stillGoing = true
			world[worldKey(x, y)] = "|"
			// look below f
			squareBelow := world[worldKey(x, y+1)]
			// if below is empty, move f down
			if squareBelow == "#" || squareBelow == "~" {
				// spread and fill
				// flow left
				var squaresLeft = 0
				var wallLeft = false
				var leftClock = 0
				var thisRow = make([]string, 0)
				for {
					leftClock++
					if leftClock > maxWidClock {
						fmt.Println("Stuck Left")
						renderWorld(world, minX, maxX, minY, maxY, frontiers)
						break
					}
					leftSquare := world[worldKey(x-squaresLeft, y)]
					leftSquareBelow := world[worldKey(x-squaresLeft, y+1)]
					if leftSquare == "#" {
						wallLeft = true
						break
					}
					key := worldKey(x-squaresLeft, y)
					world[key] = "|"
					thisRow = append(thisRow, key)
					if leftSquareBelow != "#" && leftSquareBelow != "~" {
						newFrontiers = append(newFrontiers, []int{x - squaresLeft, y + 1})
						break
					}
					squaresLeft++
				}
				// flow right
				var squaresRight = 0
				var wallRight = false
				var rightClock = 0
				for {
					rightClock++
					if rightClock > maxWidClock {
						fmt.Println("Stuck Right")
						renderWorld(world, minX, maxX, minY, maxY, frontiers)
						break
					}
					rightSquare := world[worldKey(x+squaresRight, y)]
					rightSquareBelow := world[worldKey(x+squaresRight, y+1)]
					if rightSquare == "#" {
						wallRight = true
						break
					}
					key := worldKey(x+squaresRight, y)
					world[key] = "|"
					thisRow = append(thisRow, key)
					if rightSquareBelow != "#" && rightSquareBelow != "~" {
						newFrontiers = append(newFrontiers, []int{x + squaresRight, y + 1})
						break
					}
					squaresRight++
				}
				// if blocked on both sides, fill upwards
				if wallLeft && wallRight {
					newFrontiers = append(newFrontiers, []int{x, y - 1})
					for _, key := range thisRow {
						world[key] = "~"
					}
				}
				// while bounded, fill
			} else if squareBelow == "~" {
				continue
			} else {
				// flow downwards
				newFrontiers = append(newFrontiers, []int{x, y + 1})
			}
			frontiers = newFrontiers
		}
		if !stillGoing {
			break
		}
	}
	// renderWorld(world, minX, maxX, minY, maxY, frontiers)
	fmt.Println("Part One:", waterCound(world, minX, maxX, minY, maxY, false))
	fmt.Println("Part Two:", waterCound(world, minX, maxX, minY, maxY, true))
}
