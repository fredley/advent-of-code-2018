package main

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"io/ioutil"
	"math"
	"strings"
)

type Coord struct {
	x int
	y int
}

type Square struct {
	world *[][]*Square
	x     int
	y     int
	west  bool
	north bool
	east  bool
	south bool
}

func (s *Square) PathNeighbors() []astar.Pather {
	ns := []astar.Pather{}
	for _, n := range s.neighbours() {
		ns = append(ns, n)
	}
	return ns
}

func (s *Square) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (s *Square) PathEstimatedCost(to astar.Pather) float64 {
	toSquare := to.(*Square)
	distance := math.Abs(float64(s.x-toSquare.x)) + math.Abs(float64(s.y-toSquare.y))
	return distance
}

func (s *Square) neighbours() []*Square {
	var result = make([]*Square, 0)
	if s.north {
		result = append(result, (*s.world)[s.x][s.y-1])
	}
	if s.west {
		result = append(result, (*s.world)[s.x-1][s.y])
	}
	if s.east {
		result = append(result, (*s.world)[s.x+1][s.y])
	}
	if s.south {
		result = append(result, (*s.world)[s.x][s.y+1])
	}
	return result
}

func render(world map[Coord]rune, minx int, maxx int, miny int, maxy int) string {
	result := ""
	for y := miny - 1; y <= maxy+1; y++ {
		for x := minx - 1; x < maxx+2; x++ {
			if x == 0 && y == 0 {
				result += "X"
				continue
			}
			char := world[Coord{x, y}]
			if int(char) == 0 {
				if y%2 == 0 && x%2 == 0 {
					result += fmt.Sprintf(".")
				} else {
					result += fmt.Sprintf("#")
				}
			} else {
				result += fmt.Sprintf(string(char))
			}
		}
		result += fmt.Sprintf("\n")
	}
	return result
}

func main() {
	buf, _ := ioutil.ReadFile("day-20.txt")
	input := string(buf)
	// input := "^ENWWW(NEEE|SSE(EE|N))$"
	var depth = 0
	var groups = make(map[int][]string)
	var grouppointers = make(map[int]int)
	grouppointers[0] = 0
	groups[0] = make([]string, 1)
	groups[0][0] = ""
	for _, char := range input {
		switch char {
		case '$':
			continue
		case '^':
			continue
		case 'N':
			groups[depth][grouppointers[depth]] += string(char)
		case 'S':
			groups[depth][grouppointers[depth]] += string(char)
		case 'W':
			groups[depth][grouppointers[depth]] += string(char)
		case 'E':
			groups[depth][grouppointers[depth]] += string(char)
		case '(':
			// enter group
			depth += 1
			grouppointers[depth] = 0
			groups[depth] = make([]string, 1)
			groups[depth][0] = ""
		case '|':
			grouppointers[depth] = len(groups[depth])
			groups[depth] = append(groups[depth], "")
		case ')':
			// exit group
			stringstart := groups[depth-1][grouppointers[depth-1]]
			for _, subgroup := range groups[depth] {
				groups[depth-1] = append(groups[depth-1], stringstart+subgroup)
			}
			depth -= 1
		}

	}
	world := make(map[Coord]rune)
	var minx = 0
	var miny = 0
	var maxx = 0
	var maxy = 0
	for _, path := range groups[0] {
		var x = 0
		var y = 0
		for _, char := range path {
			switch char {
			case 'N':
				world[Coord{x, y - 1}] = '-'
				y -= 2
				if miny > y {
					miny = y
				}
			case 'S':
				world[Coord{x, y + 1}] = '-'
				y += 2
				if maxy < y {
					maxy = y
				}
			case 'W':
				world[Coord{x - 1, y}] = '|'
				x -= 2
				if minx > x {
					minx = x
				}
			case 'E':
				world[Coord{x + 1, y}] = '|'
				x += 2
				if maxx < x {
					maxx = x
				}
			}
		}
	}
	worldMap := render(world, minx, maxx, miny, maxy)
	// fmt.Println(worldMap)
	lines := strings.Split(worldMap, "\n")
	worldMaze := make([][]*Square, (len(lines[0])-1)/2)
	var startx int
	var starty int
	for y, line := range lines {
		for x, char := range line {
			if x%2 == 1 && y%2 == 1 {
				newX := (x - 1) / 2
				newY := (y - 1) / 2
				if char == 'X' {
					startx = newX
					starty = newY
				}
				if y == 1 {
					worldMaze[newX] = make([]*Square, (len(lines)-1)/2)
				}
				hasNorth := (lines[y-1][x] != '#')
				hasSouth := (lines[y+1][x] != '#')
				hasWest := (lines[y][x-1] != '#')
				hasEast := (lines[y][x+1] != '#')

				worldMaze[newX][newY] = &Square{&worldMaze, newX, newY, hasWest, hasNorth, hasEast, hasSouth}
			}
		}
	}
	// fmt.Println(render(world, minx, maxx, miny, maxy))
	var mostDoors = 0.0
	var over1k = 0
	for x := 0; x < len(worldMaze[0]); x++ {
		for y := 0; y < len(worldMaze); y++ {
			_, distance, _ := astar.Path(worldMaze[startx][starty], worldMaze[x][y])
			if distance > mostDoors {
				mostDoors = distance
			}
			if distance >= 1000 {
				over1k++
			}
		}
	}
	fmt.Println("Part One:", mostDoors)
	fmt.Println("Part Two:", over1k)
}
