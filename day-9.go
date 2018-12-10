package main

import (
	"fmt"
)

func removePoint(list []int, idx int) []int {
	before := make([]int, idx)
	after := make([]int, len(list)-idx-1)
	copy(before, list[:idx])
	copy(after, list[idx+1:])
	return append(before, after...)
}

func main() {
	players := 9
	marbleWorth := 25
	var currentMarble = 0
	var circle = make([]int, 0)
	var currentPlayer = 0
	scores := make([]int, players)
	for i := 0; i < marbleWorth+1; i++ {
		if i == 0 {
			circle = append(circle, 0)
			continue
		} else if i%23 == 0 {
			removeIdx := (currentMarble - 7 + len(circle)) % len(circle)
			scores[currentPlayer] += i + circle[removeIdx]
			circle = removePoint(circle, removeIdx)
			currentMarble = removeIdx % len(circle)
		} else {
			if len(circle) == currentMarble+1 {
				circle = append(append([]int{circle[0]}, i), circle[1:]...)
				currentMarble = 1
			} else if len(circle) == currentMarble+1 {
				circle = append([]int{i}, circle...)
				currentMarble = 0
			} else {
				before := make([]int, currentMarble+2)
				after := make([]int, len(circle)-currentMarble-2)
				copy(before, circle[:currentMarble+2])
				copy(after, circle[currentMarble+2:])
				circle = append(append(before, i), after...)
				currentMarble = currentMarble + 2
			}
		}
		currentPlayer = (currentPlayer + 1) % players
	}
	var maxScore = 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}
	fmt.Println("Result", maxScore)
}
