package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"strings"
)

type Instruction struct {
	op string
	a  int
	b  int
	c  int
}

func applyOp(op string, registers *[]int, a int, b int, c int, ip int) {
	if string(op[3]) == "r" {
		b = (*registers)[b]
	}
	switch op[:3] {
	case "add":
		(*registers)[c] = (*registers)[a] + b
	case "mul":
		(*registers)[c] = (*registers)[a] * b
	case "ban":
		(*registers)[c] = (*registers)[a] & b
	case "bor":
		(*registers)[c] = (*registers)[a] | b
	case "set":
		if string(op[3]) == "r" {
			a = (*registers)[a]
		}
		(*registers)[c] = a
	case "gti":
		if a > b {
			(*registers)[c] = 1
		} else {
			(*registers)[c] = 0
		}
	case "gtr":
		if (*registers)[a] > b {
			(*registers)[c] = 1
		} else {
			(*registers)[c] = 0
		}
	case "eqi":
		if a == b {
			(*registers)[c] = 1
		} else {
			(*registers)[c] = 0
		}
	case "eqr":
		if (*registers)[a] == b {
			(*registers)[c] = 1
		} else {
			(*registers)[c] = 0
		}
	}
	(*registers)[ip]++
}

func main() {
	buf, _ := ioutil.ReadFile("day-19.txt")
	input := string(buf)
	lines := strings.Split(input, "\n")
	instructions := make([]Instruction, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, " ")
		a, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])
		c, _ := strconv.Atoi(parts[3])
		instructions[i] = Instruction{parts[0], a, b, c}
	}
	registers := []int{0, 0, 0, 0, 0, 0}
	instructionRegister := 1
	for {
		ptr := registers[instructionRegister]
		if ptr < 0 || ptr >= len(instructions) {
			fmt.Println("Part One:", registers[0])
			break
		}
		i := instructions[ptr]
		applyOp(i.op, &registers, i.a, i.b, i.c, instructionRegister)
	}
	// Part 2
	n := 10551408
	var res = 0
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			res += i
		}
	}
	fmt.Println("Part Two:", res)
}
