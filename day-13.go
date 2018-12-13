package main

import (
	"fmt"
	"sort"
	"strings"
)

type Cart struct {
	x           int
	y           int
	dir         int
	isCrashed   bool
	turnCounter int
	collided    bool
}

type Carts []Cart

func (s Carts) Len() int {
	return len(s)
}
func (s Carts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Carts) Less(i, j int) bool {
	l1 := s[i]
	l2 := s[j]
	if l1.y < l2.y || (l1.y == l2.y && l1.x < l2.x) {
		return true
	}
	return false
}

func checkCollisions(cart Cart, carts []Cart) (bool, []Cart) {
	var collided = false
	var cartOneFound = false
	var cart1 = 0
	var cart2 = 0
	for i, c := range carts {
		if !c.collided && c.x == cart.x && c.y == cart.y {
			if !cartOneFound {
				cartOneFound = true
				cart1 = i
			} else {
				cart2 = i
				collided = true
			}
		}
	}
	if collided {
		carts[cart1].collided = true
		carts[cart2].collided = true
	}
	return collided, carts
}

func render(grid []string, carts []Cart) string {
	var result = ""
	for j, line := range grid {
		for i, char := range line {
			var found = false
			for _, cart := range carts {
				if cart.x == i && cart.y == j && !cart.collided {
					if cart.dir == 0 {
						result += "<"
					} else if cart.dir == 1 {
						result += ">"
					} else if cart.dir == 2 {
						result += "^"
					} else if cart.dir == 3 {
						result += "v"
					}
					found = true
					break
				}
			}
			if !found {
				result += string(char)
			}
		}
		result += "\n"
	}
	return result
}

func main() {
	const input = ``

	lines := strings.Split(input, "\n")
	// clean map
	var cleanMap = ""
	var carts = make([]Cart, 0)
	for j, line := range lines {
		for i, char := range line {
			var strChar = string(char)
			if strChar == "<" {
				carts = append(carts, Cart{i, j, 0, false, 0, false})
				strChar = "-"
			} else if strChar == ">" {
				carts = append(carts, Cart{i, j, 1, false, 0, false})
				strChar = "-"
			} else if strChar == "^" {
				carts = append(carts, Cart{i, j, 2, false, 0, false})
				strChar = "|"
			} else if strChar == "v" {
				carts = append(carts, Cart{i, j, 3, false, 0, false})
				strChar = "|"
			}
			cleanMap += strChar
		}
		cleanMap += "\n"
	}
	var firstCollision = false
	for {
		sort.Sort(Carts(carts))
		for i := 0; i < len(carts); i++ {
			cart := carts[i]
			// check if collided
			if cart.collided {
				continue
			}
			// move cart
			if cart.dir <= 1 {
				carts[i].x += cart.dir*2 - 1
			} else {
				carts[i].y += (cart.dir-2)*2 - 1
			}
			//check for collision
			var hasCollided bool
			hasCollided, carts = checkCollisions(carts[i], carts)
			if hasCollided {
				if !firstCollision {
					fmt.Printf("Part One: %d,%d\n", cart.x, cart.y)
					firstCollision = true
				}
				continue
			}
			// check current square and update direction accordingly
			newSquare := string(lines[carts[i].y][carts[i].x])
			if newSquare == "\\" {
				carts[i].dir = (cart.dir + 2) % 4
			} else if newSquare == "/" {
				if cart.dir%2 == 0 {
					carts[i].dir = (cart.dir + 3) % 4
				} else {
					carts[i].dir = (cart.dir + 1) % 4
				}
			} else if newSquare == "+" {
				if cart.turnCounter == 0 {
					// turn left
					if cart.dir <= 1 {
						carts[i].dir = (cart.dir + (4 + (cart.dir*2 - 1))) % 4
					} else {
						carts[i].dir = (cart.dir + 2) % 4
					}
				} else if cart.turnCounter == 2 {
					// turn right
					if cart.dir <= 1 {
						carts[i].dir = cart.dir + 2
					} else {
						carts[i].dir = (cart.dir + (4 + ((cart.dir-2)*2 - 1))) % 4
					}
				}
				carts[i].turnCounter = (cart.turnCounter + 1) % 3
			}
		}
		var alive = 0
		var aliveIndex = 0
		for i := 0; i < len(carts); i++ {
			if !carts[i].collided {
				alive += 1
				aliveIndex = i
			}
		}
		if alive == 1 {
			fmt.Printf("Part Two: %d,%d\n", carts[aliveIndex].x, carts[aliveIndex].y)
			break
		}
	}

}
