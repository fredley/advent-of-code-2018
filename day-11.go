package main

import (
	"fmt"
)

const SERIAL_NUMBER = 8
const GRID_SIZE = 300

type Cell struct {
	x        int
	y        int
	power    int
	powerSet bool
}

func (c *Cell) getPower() int {
	if !c.powerSet {
		rackID := c.x + 10
		c.power = ((((rackID * c.y) + SERIAL_NUMBER) * rackID) / 100 % 10) - 5
		c.powerSet = true
	}
	return c.power
}

func gridPower(grid [][]Cell, x int, y int, squareSize int) int {
	var power = 0
	for i := 0; i <= squareSize-1; i++ {
		for j := 0; j <= squareSize-1; j++ {
			power += grid[x+i][y+j].getPower()
		}
	}
	return power
}

func maxGridPower(grid [][]Cell, squareSize int) (int, int, int) {
	var maxX = 0
	var maxY = 0
	var maxPower = -9
	for i := 0; i <= GRID_SIZE-squareSize; i++ {
		for j := 0; j <= GRID_SIZE-squareSize; j++ {
			power := gridPower(grid, i, j, squareSize)
			if power > maxPower {
				maxPower = power
				maxX = i + 1
				maxY = j + 1
			}
		}
	}
	return maxPower, maxX, maxY
}

func main() {
	grid := make([][]Cell, GRID_SIZE)
	for i := 0; i < GRID_SIZE; i++ {
		grid[i] = make([]Cell, GRID_SIZE)
		for j := 0; j < GRID_SIZE; j++ {
			grid[i][j] = Cell{i + 1, j + 1, 0, false}
		}
	}
	var maxPower = -9
	var maxSquareSize = 0
	var maxX = 0
	var maxY = 0
	for s := 3; s <= 300; s++ {
		power, x, y := maxGridPower(grid, s)
		if s == 3 {
			fmt.Printf("Part one: %d,%d\n", x, y)
		}
		if power > maxPower {
			maxPower = power
			maxX = x
			maxY = y
			maxSquareSize = s
		}
	}
	fmt.Printf("Part two: %d,%d,%d\n", maxX, maxY, maxSquareSize)

}
