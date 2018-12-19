package main

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

type Square struct {
	x      int
	y      int
	isWall bool
	unit   *Unit
	world  *World
}

type Unit struct {
	isElf   bool
	isAlive bool
	hp      int
	ap      int
	square  *Square
}

func (u Unit) enemies(otherUnits []*Unit) []*Unit {
	var result = make([]*Unit, 0)
	for _, unit := range otherUnits {
		if (unit.isElf != u.isElf) && unit.isAlive {
			result = append(result, unit)
		}
	}
	return result
}

func allTargets(units []*Unit) []*Square {
	var result = make([]*Square, 0)
	for _, unit := range units {
		result = append(result, unit.square.neighbours()...)
	}
	return result
}

func (s *Square) free() bool {
	return !s.isWall && (s.unit == nil || !s.unit.isAlive)
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
	// fmt.Println("estimate", s.x, s.y, "to", toSquare.x, toSquare.y)
	distance := math.Abs(float64(s.x-toSquare.x)) + math.Abs(float64(s.y-toSquare.y))
	return distance
}

func (s *Square) neighbours() []*Square {
	var result = make([]*Square, 0)
	// TODO BOUNDS
	north := &s.world.squares[s.x][s.y-1]
	west := &s.world.squares[s.x-1][s.y]
	east := &s.world.squares[s.x+1][s.y]
	south := &s.world.squares[s.x][s.y+1]
	if north.free() {
		result = append(result, north)
	}
	if west.free() {
		result = append(result, west)
	}
	if east.free() {
		result = append(result, east)
	}
	if south.free() {
		result = append(result, south)
	}
	return result
}

func (u Unit) canAttack(otherUnit *Unit) bool {
	if u.square.x == otherUnit.square.x {
		if u.square.y == otherUnit.square.y+1 || u.square.y == otherUnit.square.y-1 {
			return true
		}
	} else if u.square.y == otherUnit.square.y {
		if u.square.x == otherUnit.square.x+1 || u.square.x == otherUnit.square.x-1 {
			return true
		}
	}
	return false
}
func (u Unit) name() string {
	if u.isElf {
		return "Elf"
	}
	return "Goblin"
}

func filterUnique(squares []*Square) []*Square {
	sort.Sort(Squares(squares))
	var result = make([]*Square, 0)
	var prev *Square
	for i, s := range squares {
		if i == 0 {
			prev = s
			result = append(result, s)
			continue
		}
		if (s.x == prev.x && s.y == prev.y) || !s.free() {
			continue
		}
		prev = s
		result = append(result, s)
	}
	return result
}

type Units []*Unit

func (s Units) Len() int {
	return len(s)
}
func (s Units) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Units) Less(i, j int) bool {
	l1 := s[i]
	l2 := s[j]
	if l1.square.y < l2.square.y || (l1.square.y == l2.square.y && l1.square.x < l2.square.x) {
		return true
	}
	return false
}

type Squares []*Square

func (s Squares) Len() int {
	return len(s)
}
func (s Squares) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Squares) Less(i, j int) bool {
	l1 := s[i]
	l2 := s[j]
	if l1.y < l2.y || (l1.y == l2.y && l1.x < l2.x) {
		return true
	}
	return false
}

type World struct {
	units   []*Unit
	squares [][]Square
	width   int
	height  int
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func renderWorld(world World) {
	for y := 0; y < world.height; y++ {
		rowAfter := " "
		for x := 0; x < world.width; x++ {
			square := world.squares[x][y]
			if square.isWall {
				fmt.Printf("#")
			} else if square.unit != nil && square.unit.isAlive {
				if len(rowAfter) > 1 {
					rowAfter += ", "
				}
				if square.unit.isElf {
					fmt.Printf("E")
					rowAfter += "E"
				} else {
					fmt.Printf("G")
					rowAfter += "G"
				}
				rowAfter += "(" + fmt.Sprintf("%d", square.unit.hp) + ")"
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf(rowAfter + "\n")
	}
}

func (world *World) iterate() bool {
	sort.Sort(Units(world.units))
	for idx, unit := range world.units {
		if !unit.isAlive {
			continue
		}
		enemies := unit.enemies(world.units)
		if len(enemies) == 0 {
			return false
		}
		var shouldMove = true
		for _, enemy := range enemies {
			xdiff := math.Abs(float64(enemy.square.x - unit.square.x))
			ydiff := math.Abs(float64(enemy.square.y - unit.square.y))
			if xdiff+ydiff == 1 {
				shouldMove = false
				break
			}
		}
		targets := filterUnique(allTargets(enemies))
		if len(targets) == 0 {
			shouldMove = false
		}
		if shouldMove {
			// Find nearest reachable neighbours
			// var minEnemies = make([]Unit, 0)
			var minSquares = make([]*Square, 0)
			var minDistance = 999.9

			// fmt.Println("\n", unit.name(), "with", len(targets), "targets from", len(enemies), "enemies")
			for _, target := range targets {
				_, distance, found := astar.Path(unit.square, target)
				if !found {
					continue
				}
				if distance == minDistance {
					minDistance = distance
					// minEnemies = append(minEnemies, enemy)
					minSquares = append(minSquares, target)
				} else if distance < minDistance {
					minDistance = distance
					// minEnemies = []Unit{enemy}
					minSquares = []*Square{target}
				}
			}
			if len(minSquares) > 1 {
				sort.Sort(Squares(minSquares))
			} else {
				// fmt.Println("Next step is", path[0])
			}
			if len(minSquares) > 0 {
				minSteps := 999.9
				validSteps := make([]*Square, 0)
				for _, n := range unit.square.neighbours() {
					_, steps, found := astar.Path(n, minSquares[0])
					if !found {
						continue
					}
					if steps == minSteps {
						validSteps = append(validSteps, n)
					} else if steps < minSteps {
						validSteps = []*Square{n}
						minSteps = steps
					}
				}
				if len(validSteps) > 1 {
					sort.Sort(Squares(validSteps))
				}
				// fmt.Println("Next step is from", unit.square.x, unit.square.y, "to", validSteps[0].x, validSteps[0].y)
				world.squares[unit.square.x][unit.square.y].unit = nil
				validSteps[0].unit = world.units[idx]
				world.units[idx].square = validSteps[0]
			} else {
				// fmt.Println("No valid moves")
			}

		} else {
			// fmt.Println(unit.name(), "should not move from", unit.square.x, unit.square.y)
		}
		possibleEnemies := make([]*Unit, 0)
		var minHp = 201
		for _, enemy := range enemies {
			if !unit.canAttack(enemy) {
				continue
			}
			if enemy.hp < minHp {
				possibleEnemies = []*Unit{enemy}
				minHp = enemy.hp
			} else if enemy.hp == minHp {
				possibleEnemies = append(possibleEnemies, enemy)
				minHp = enemy.hp
			}
		}
		if len(possibleEnemies) > 0 {
			sort.Sort(Units(possibleEnemies))
			chosenEnemy := possibleEnemies[0]
			// fmt.Println("Attacking", chosenEnemy)
			chosenEnemy.hp -= unit.ap
			if chosenEnemy.hp <= 0 {
				chosenEnemy.isAlive = false
			}
		}
	}
	return true
}

func testWorld(lines []string, attackPower int) int {
	units := make([]*Unit, 0)
	squares := make([][]Square, len(lines[0]))
	world := World{units, squares, len(lines[0]), len(lines)}
	numGoblins := 0
	numELves := 0
	for y, line := range lines {
		for x, char := range line {
			if y == 0 {
				squares[x] = make([]Square, len(lines))
			}
			var unit Unit
			if char != '#' && char != '.' {
				var power int
				if char == 'E' {
					numELves++
					power = attackPower
				} else {
					numGoblins++
					power = 3
				}
				unit = Unit{char == 'E', true, 200, power, nil}
			}
			square := Square{x, y, char == '#', &unit, &world}
			world.squares[x][y] = square
			unit.square = &square
			if unit.isAlive {
				world.units = append(world.units, &unit)
			}
		}
	}
	var stillGoing = true
	for i := 0; i < 10000; i++ {
		stillGoing = world.iterate()
		if !stillGoing {
			var score = 0
			var elvesAlive = 0
			for _, unit := range world.units {
				if unit.isAlive {
					score += unit.hp
					if unit.isElf {
						elvesAlive++
					}
				}
			}
			if attackPower == 3 {
				fmt.Println("Part One:", i*score)
			}
			if elvesAlive == numELves {
				return i * score
			} else {
				return -1
			}
		}
	}
	fmt.Println("timed out")
	return -1
}

func main() {
	buf, _ := ioutil.ReadFile("day-15.txt")
	lines := strings.Split(string(buf), "\n")
	for i := 3; i < 999; i++ {
		result := testWorld(lines, i)
		if result > 0 {
			fmt.Println("Part Two:", result)
			break
		}
	}
}
