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
	buf, _ := ioutil.ReadFile("day-21.txt")
	input := string(buf)
	lines := strings.Split(input, "\n")
	instructions := make([]Instruction, len(lines)-1)
	var instructionRegister = 0
	for i, line := range lines {
		if i == 0 {
			instructionRegister, _ = strconv.Atoi(strings.Split(line, " ")[1])
			continue
		}
		parts := strings.Split(line, " ")
		a, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])
		c, _ := strconv.Atoi(parts[3])
		instructions[i-1] = Instruction{parts[0], a, b, c}
	}
	var registerOne = 0
	seenR5 := make(map[int]bool)
	for {
		registers := []int{registerOne, 0, 0, 0, 0, 0}
		seenStates := make(map[string]bool)
		lowest := 999999999
		cycles := 0
		for {
			ptr := registers[instructionRegister]
			if ptr < 0 || ptr >= len(instructions) {
				fmt.Println("Complete in", cycles, "at", registerOne)
				if lowest > cycles {
					lowest = cycles
					fmt.Println("New lowest!", lowest, "with register one at", registerOne)
				}
				break
			}
			i := instructions[ptr]
			if ptr == 28 {
				if !seenR5[registers[5]] {
					fmt.Println("Target value", registers[5])
				}
				seenR5[registers[5]] = true
			}
			applyOp(i.op, &registers, i.a, i.b, i.c, instructionRegister)
			stateKey := string(registers[0]) + string(registers[1]) + string(registers[2]) + string(registers[3]) + string(registers[4]) + string(registers[5])
			if seenStates[stateKey] {
				break
			}
			seenStates[stateKey] = true
			cycles++
		}
		registerOne++
	}
}
