package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func sequentialize(sequence []int, target []int, newEl int) ([]int, bool) {
	if len(target) == len(sequence) {
		return sequence, false
	}
	if target[len(sequence)] == newEl {
		sequence = append(sequence, newEl)
		if len(sequence) == len(target) {
			return sequence, true
		}
		return sequence, false
	} else if target[0] == newEl {
		return []int{newEl}, false
	}
	return make([]int, 0), false
}

func iterate(recipes []int, elf1 int, elf2 int, sequence []int, target []int) ([]int, int, int, []int, bool) {
	newSum := recipes[elf1] + recipes[elf2]
	var done = false
	if newSum < 10 {
		recipes = append(recipes, newSum)
		sequence, done = sequentialize(sequence, target, newSum)
		if done {
			return recipes, elf1, elf2, sequence, true
		}
	} else {
		recipes = append(append(recipes, 1), newSum-10)
		sequence, done = sequentialize(sequence, target, 1)
		if done {
			return recipes, elf1, elf2, sequence, true
		}
		sequence, done = sequentialize(sequence, target, newSum-10)
		if done {
			return recipes, elf1, elf2, sequence, true
		}
	}
	elf1 = (elf1 + recipes[elf1] + 1) % len(recipes)
	elf2 = (elf2 + recipes[elf2] + 1) % len(recipes)
	return recipes, elf1, elf2, sequence, false
}

func render(recipes []int, elf1 int, elf2 int) {
	for i := 0; i < len(recipes); i++ {
		if i == elf1 {
			fmt.Printf("(%d)", recipes[i])
		} else if i == elf2 {
			fmt.Printf("[%d]", recipes[i])
		} else {
			fmt.Printf(" %d ", recipes[i])
		}
	}
	fmt.Printf("\n")
}

func renderScore(slice []int) string {
	var result = ""
	for i := 0; i < len(slice); i++ {
		result += strconv.Itoa(slice[i])
	}
	return result
}

func main() {
	var recipes = make([]int, 2)
	recipes[0] = 3
	recipes[1] = 7

	var elf1 = 0
	var elf2 = 1

	target := 2018
	targetSequence := []int{5, 9, 4, 1, 4}

	var sequence = make([]int, 0)
	var done = false
	var done1 = false
	var done2 = false

	for {
		recipes, elf1, elf2, sequence, done = iterate(recipes, elf1, elf2, sequence, targetSequence)
		if len(recipes) > target+10 && !done1 {
			fmt.Println("Part One:", renderScore(recipes[target:target+10]))
			done1 = true
		}
		if done && !done2 {
			var resultseq = recipes[len(recipes)-len(sequence)-1:]
			var partTwo = len(recipes) - len(sequence)
			if reflect.DeepEqual(resultseq[:len(resultseq)-1], targetSequence) {
				partTwo--
			}
			fmt.Println("Part Two:", partTwo)
			done2 = true
		}
		if done1 && done2 {
			return
		}
	}

}
