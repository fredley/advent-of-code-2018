package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func applyOp(op string, registers *[]int, a int, b int, c int) {
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
}

func isEqual(r1 []int, r2 []int) bool {
	for i, val := range r1 {
		if val != r2[i] {
			return false
		}
	}
	return true
}

func copyRegisters(input []int) []int {
	output := make([]int, len(input))
	for i := range input {
		output[i] = input[i]
	}
	return output
}

func main() {
	buf, _ := ioutil.ReadFile("day-16.txt")
	input := string(buf)
	lines := strings.Split(input, "\n")
	var countOps = 0
	var canbe = make(map[int]map[string]bool)
	var ops = []string{"addr", "addi", "mulr", "muli", "banr", "bani", "borr", "bori", "setr", "seti", "gtir", "gtri", "gtrr", "eqir", "eqri", "eqrr"}
	for i := 0; i < 16; i++ {
		canbe[i] = make(map[string]bool)
		for _, op := range ops {
			canbe[i][op] = true
		}
	}
	for i := 0; i < len(lines); i++ {
		inputline := lines[i]
		if len(inputline) == 0 {
			continue
		}
		if inputline[:6] == "Before" {
			sampleInput := strings.Split(strings.Split(inputline, "[")[1], "]")[0]
			sampleInputRegisters := make([]int, 4)
			inputValues := strings.Split(sampleInput, " ")
			for j := range sampleInputRegisters {
				value, _ := strconv.ParseInt(strings.Split(inputValues[j], ",")[0], 10, 0)
				sampleInputRegisters[j] = int(value)
			}
			i++
			opline := lines[i]
			opValues := strings.Split(opline, " ")
			var opcode int
			var a int
			var b int
			var c int
			for j := 0; j < len(opValues); j++ {
				value, _ := strconv.ParseInt(opValues[j], 10, 0)
				if j == 0 {
					opcode = int(value)
				} else if j == 1 {
					a = int(value)
				} else if j == 2 {
					b = int(value)
				} else if j == 3 {
					c = int(value)
				}
			}
			i++
			outputline := lines[i]
			sampleOutput := strings.Split(strings.Split(outputline, "[")[1], "]")[0]
			sampleOutputRegisters := make([]int, 4)
			outputValues := strings.Split(sampleOutput, " ")
			for j := range sampleOutputRegisters {
				value, _ := strconv.ParseInt(strings.Split(outputValues[j], ",")[0], 10, 0)
				sampleOutputRegisters[j] = int(value)
			}
			// Actually test
			var countSame = 0
			for _, op := range ops {
				var testRegisters = copyRegisters(sampleInputRegisters)
				applyOp(op, &testRegisters, a, b, c)
				if isEqual(testRegisters, sampleOutputRegisters) {
					countSame++
				} else {
					canbe[opcode][op] = false
				}
			}
			if countSame >= 3 {
				countOps++
			}
		}
	}
	fmt.Println("Part One:", countOps)
	var opcodes = make(map[int]string)
	var opsfound = 0
	var found = make(map[string]bool)
	for _, opname := range ops {
		found[opname] = false
	}
	for {
		var opLocations = make(map[string][]int)
		for _, op := range ops {
			opLocations[op] = make([]int, 0)
		}
		for opcode, canbes := range canbe {
			var numOpts = 0
			var optCandidate string
			for opname, couldbe := range canbes {
				if !found[opname] && couldbe {
					optCandidate = opname
					numOpts++
					opLocations[opname] = append(opLocations[opname], opcode)
				}
			}
			if numOpts == 1 {
				found[optCandidate] = true
				opcodes[opcode] = optCandidate
				opsfound++
				fmt.Println("Found", optCandidate, opcode)
			}
		}
		for op, locations := range opLocations {
			if found[op] {
				continue
			}
			if len(locations) == 1 {
				found[op] = true
				opcodes[locations[0]] = op
				opsfound++
				fmt.Println("Found", op, locations[0])
			}
		}
		if opsfound == 16 {
			break
		}
	}
	buf2, _ := ioutil.ReadFile("day-16-2.txt")
	input2 := string(buf2)
	programLines := strings.Split(input2, "\n")
	registers := []int{0, 0, 0, 0}
	for _, line := range programLines {
		if len(line) == 0 {
			break
		}
		split := strings.Split(line, " ")
		opval, _ := strconv.ParseInt(split[0], 10, 0)
		aval, _ := strconv.ParseInt(split[1], 10, 0)
		bval, _ := strconv.ParseInt(split[2], 10, 0)
		cval, _ := strconv.ParseInt(split[3], 10, 0)
		applyOp(opcodes[int(opval)], &registers, int(aval), int(bval), int(cval))
	}
	fmt.Println("Part Two:", registers[0])
}
