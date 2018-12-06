package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Coord struct {
	x int
	y int
}

type Grid struct {
	grid   map[Coord]int
	width  int
	height int
}

func printGrid(grid Grid) {
	for i := 0; i < grid.width; i++ {
		for j := 0; j < grid.height; j++ {
			var c = grid.grid[Coord{i, j}]
			if j == 0 {
				fmt.Print("\n")
			}
			fmt.Print(c)
		}
	}
	fmt.Print("\n")
}

func main() {
	const input = ``
	var lines = strings.Split(input, "\n")
	// parse coords, find max size
	var coords = make([]Coord, len(lines))
	var maxX float64
	var maxY float64
	for i, line := range lines {
		var split = strings.Split(line, ", ")
		var y, _ = strconv.ParseInt(split[0], 10, 0)
		var x, _ = strconv.ParseInt(split[1], 10, 0)
		maxX = math.Max(float64(x)+1, maxX)
		maxY = math.Max(float64(y)+2, maxY)
		coords[i] = Coord{int(x), int(y)}
	}
	// create a grid
	var grid = Grid{make(map[Coord]int), int(maxX), int(maxY)}
	for i := 0; i < grid.width; i++ {
		for j := 0; j < grid.height; j++ {
			grid.grid[Coord{i, j}] = 0
		}
	}
	var id = 1
	for _, coord := range coords {
		grid.grid[coord] = id
		id++
	}
	var counts = make(map[int]int)
	var edges = make(map[int]int)
	var safe = 0
	// find nearest for each space
	for i := 0; i < grid.width; i++ {
		for j := 0; j < grid.height; j++ {
			var sumDistance = 0
			var minDistance = maxX + maxY
			var minChar = 0
			var isDouble = false
			var isThis = false
			for _, coord := range coords {
				var distance = math.Abs(float64(coord.x-i)) + math.Abs(float64(coord.y-j))
				sumDistance += int(distance)
				if distance == 0 {
					isThis = true
					counts[grid.grid[coord]]++
				}
				if distance == minDistance {
					isDouble = true
				} else if distance < minDistance {
					isDouble = false
					minDistance = distance
					minChar = grid.grid[coord]
				}
			}
			if sumDistance < 10000 {
				safe++
			}
			if isDouble || isThis {
				continue
			}
			grid.grid[Coord{i, j}] = minChar
			counts[minChar]++
			if i == 0 || j == 0 || i == grid.width-1 || j == grid.height-1 {
				// This is on the edge
				edges[minChar] = 1
			}
		}
	}
	// find largest area
	var maxArea = 0
	for n, count := range counts {
		// ignore any areas that have points on the boundary
		if edges[n] == 0 {
			if count > maxArea {
				maxArea = count
			}
		}
	}
	fmt.Println("Part One:", maxArea)
	fmt.Println("Part Two:", safe)
}
