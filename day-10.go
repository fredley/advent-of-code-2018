package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	x  int
	y  int
	dx int
	dy int
}

func getDims(points []Point) (int, int, int, int) {
	var minX = 0
	var minY = 0
	var maxX = 0
	var maxY = 0
	for _, point := range points {
		if point.x < minX {
			minX = point.x
		}
		if point.x > maxX {
			maxX = point.x
		}
		if point.y < minY {
			minY = point.y
		}
		if point.y > maxY {
			maxY = point.y
		}
	}
	width := maxX - minX
	height := maxY - minY
	return width, height, minX, minY
}

func drawPoints(points []Point) {
	width, height, minX, minY := getDims(points)
	grid := make([][]bool, width+1)
	for _, point := range points {
		if grid[point.x-minX] == nil {
			grid[point.x-minX] = make([]bool, height+1)
		}
		grid[point.x-minX][point.y-minY] = true
	}
	var res = ""
	for j := 0; j < height+1; j++ {
		for i := 0; i < width+1; i++ {
			if grid[i] != nil && grid[i][j] {
				res += "#"
			} else {
				res += " "
			}
		}
		res += "\n"
	}
	fmt.Println(res)
}

func iterate(points []Point) []Point {
	for i, point := range points {
		points[i].x += point.dx
		points[i].y += point.dy
	}
	return points
}
func reverse(points []Point) []Point {
	for i, point := range points {
		points[i].x -= point.dx
		points[i].y -= point.dy
	}
	return points
}

func main() {
	const input = ``
	split := strings.Split(input, "\n")
	points := make([]Point, len(split))
	for i, line := range split {
		x, _ := strconv.ParseInt(strings.Trim(strings.Split(strings.Split(line, "<")[1], ",")[0], " "), 10, 0)
		y, _ := strconv.ParseInt(strings.Trim(strings.Split(strings.Split(line, ",")[1], ">")[0], " "), 10, 0)
		dx, _ := strconv.ParseInt(strings.Trim(strings.Split(strings.Split(line, "<")[2], ",")[0], " "), 10, 0)
		dy, _ := strconv.ParseInt(strings.Trim(strings.Split(strings.Split(line, ",")[2], ">")[0], " "), 10, 0)
		points[i] = Point{int(x), int(y), int(dx), int(dy)}
	}
	var minHeight = 1000000
	var i = 0
	for {
		points = iterate(points)
		_, height, _, _ := getDims(points)
		if height < minHeight {
			minHeight = height
			i += 1
		} else {
			points = reverse(points)
			drawPoints(points)
			fmt.Println("Part two:", i)
			break
		}
	}
}
