package main

import (
	"fmt"
	"github.com/RyanCarrier/dijkstra"
)

type Square struct {
	x        int
	y        int
	z        int
	passable bool
	id       int
	world    *[][][]*Square
}

func (s *Square) neighbours() []*Square {
	var result = make([]*Square, 0)
	// TODO BOUNDS
	if s.x != 0 {
		west := (*s.world)[s.x-1][s.y][s.z]
		if west.passable {
			result = append(result, west)
		}
	}
	if s.y != 0 {
		north := (*s.world)[s.x][s.y-1][s.z]
		if north.passable {
			result = append(result, north)
		}
	}
	if s.x != len(*s.world)-1 {
		east := (*s.world)[s.x+1][s.y][s.z]
		if east.passable {
			result = append(result, east)
		}
	}
	if s.y != len((*s.world)[0])-1 {
		south := (*s.world)[s.x][s.y+1][s.z]
		if south.passable {
			result = append(result, south)
		}
	}
	for i := 0; i < 3; i++ {
		toolChange := (*s.world)[s.x][s.y][i]
		if i != s.z && toolChange.passable {
			result = append(result, toolChange)
		}
	}
	return result
}

const DEPTH = 8787

var idxmemset = make(map[int]bool)
var idxmemo = make(map[int]int)

func geologicIndex(x int, y int) int {
	key := 9999*x + y
	if idxmemset[key] {
		return idxmemo[key]
	}
	if x == 0 {
		if y == 0 {
			res := 0
			idxmemo[key] = res
			idxmemset[key] = true
			return res
		}
		res := y * 48271
		idxmemo[key] = res
		idxmemset[key] = true
		return res
	}
	if y == 0 {
		res := x * 16807
		idxmemo[key] = res
		idxmemset[key] = true
		return res
	}
	res := erosionLevel(x-1, y) * erosionLevel(x, y-1)
	idxmemo[key] = res
	idxmemset[key] = true
	return res
}

var geomemset = make(map[int]bool)
var geomemo = make(map[int]int)

func erosionLevel(x int, y int) int {
	key := 9999*x + y
	if geomemset[key] {
		return geomemo[key]
	}
	res := (geologicIndex(x, y) + DEPTH) % 20183
	geomemo[key] = res
	geomemset[key] = true
	return res
}

func riskLevel(x int, y int) uint {
	return uint(erosionLevel(x, y) % 3)
}

func renderWorld(world [][]int, targetX int, targetY int) {
	chars := ".=|"
	for y := 0; y < len(world[0]); y++ {
		for x := 0; x < len(world); x++ {
			if (x == 0 && y == 0) || (x == targetX-1 && y == targetY-1) {
				fmt.Printf("X")
			} else {
				fmt.Printf(string(chars[world[x][y]]))
			}
		}
		fmt.Printf("\n")
	}
}

func renderTool(world [][]int, targetX int, targetY int, tool int) {
	// torch 0, gear 1, neither 2
	// rocky 0, wet  1, narrow  2
	chars := "O OOO  OO"[tool*3 : tool*3+3]
	for y := 0; y < len(world[0]); y++ {
		for x := 0; x < len(world); x++ {
			if (x == 0 && y == 0) || (x == targetX-1 && y == targetY-1) {
				fmt.Printf("X")
			} else {
				fmt.Printf(string(chars[world[x][y]]))
			}
		}
		fmt.Printf("\n")
	}
}

func calculateRisk(world [][]int, targetX int, targetY int) int {
	res := 0
	for y := 0; y < len(world[0]); y++ {
		for x := 0; x < len(world); x++ {
			if (x == 0 && y == 0) || (x == targetX && y == targetY) {
				continue
			}
			res += world[x][y]
		}
	}
	return res
}

func renderLayers(world [][][]*Square, targetX int, targetY int, paths map[int]bool) {
	for y := 0; y < len(world[0]); y++ {
		for z := 0; z < 3; z++ {
			for x := 0; x < len(world); x++ {
				if world[x][y][z].passable {
					if x == targetX && y == targetY {
						fmt.Printf("X")
					} else if paths[x*5+y*3+z] {
						fmt.Printf("+")
					} else {
						fmt.Printf(" ")
					}
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf(" - ")
		}
		fmt.Printf("\n")
	}
}

func main() {
	targetX := 10 + 1
	// targetY := 725
	targetY := 725 + 1
	world := make([][]int, targetX)
	for i := range world {
		world[i] = make([]int, targetY)
		for j := range world[i] {
			if i == targetX-1 && j == targetY-1 {
				world[i][j] = (DEPTH % 20183) % 3
			} else {
				world[i][j] = int(riskLevel(i, j))
			}
		}
	}
	fmt.Println("Part One:", calculateRisk(world, targetX, targetY))

	const gear = 1<<2 + 1<<1
	const neither = 1<<1 + 1
	const torch = 1<<2 + 1

	tools := []int{gear, neither, torch}

	offset := 250
	solveWorld := make([][][]*Square, targetX+offset)
	squareList := make([]*Square, (targetX+offset)*(targetY+offset)*3)
	graph := dijkstra.NewGraph()
	idx := 0
	for i := 0; i < targetX+offset; i++ {
		solveWorld[i] = make([][]*Square, targetY+offset)
		for j := 0; j < targetY+offset; j++ {
			solveWorld[i][j] = make([]*Square, 3)
			for z := 0; z < 3; z++ {
				if i == targetX-1 && j == targetY-1 {
					s := Square{i, j, z, (tools[z] & (1 << ((DEPTH % 20183) % 3))) > 0, idx, &solveWorld}
					solveWorld[i][j][z] = &s
					squareList[idx] = &s
				} else {
					s := Square{i, j, z, (tools[z] & (1 << riskLevel(i, j))) > 0, idx, &solveWorld}
					solveWorld[i][j][z] = &s
					squareList[idx] = &s
				}
				graph.AddVertex(idx)
				idx++
			}
		}
	}
	// renderLayers(solveWorld, 0, 0)

	start := 0
	end := 0
	var distance int64
	for x := 0; x < len(solveWorld); x++ {
		for y := 0; y < len(solveWorld[0]); y++ {
			for z := 0; z < 3; z++ {
				s := solveWorld[x][y][z]
				if x == 0 && y == 0 && z == 2 {
					start = s.id
				} else if x == targetX-1 && y == targetY-1 && z == 2 {
					end = s.id
				}
				for _, n := range s.neighbours() {
					if s.z == n.z {
						distance = 1
					} else {
						distance = 7
					}
					graph.AddArc(s.id, n.id, distance)
				}
			}
		}
	}

	pathLength, err := graph.Shortest(start, end)
	if err != nil {
		fmt.Println(err)
	}
	paths := make(map[int]bool)
	for _, idx := range pathLength.Path {
		s := squareList[idx]
		paths[s.x*5+s.y*3+s.z] = true
	}
	// renderLayers(solveWorld, targetX, targetY, paths)
	fmt.Println("Part Two:", pathLength.Distance)
}
