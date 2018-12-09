package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	children []Node
	metadata []int
}

func parseNode(input []int) (Node, []int) {
	numChildren := input[0]
	numMetadata := input[1]
	var children = make([]Node, 0)
	var remainder = input[2:]
	for i := 0; i < numChildren; i++ {
		child, rem := parseNode(remainder)
		remainder = rem
		children = append(children, child)
	}
	return Node{children, remainder[:numMetadata]}, remainder[numMetadata:]
}

func sumMeta(node Node) int {
	var sum = 0
	for _, child := range node.children {
		sum += sumMeta(child)
	}
	for _, val := range node.metadata {
		sum += val
	}
	return sum
}

func sumRoot(node Node) int {
	var sum = 0
	if len(node.children) == 0 {
		for _, val := range node.metadata {
			sum += val
		}
	} else {
		for _, val := range node.metadata {
			if len(node.children) > val-1 {
				sum += sumRoot(node.children[val-1])
			}
		}
	}
	return sum
}

func main() {
	const input = ``
	split := strings.Split(input, " ")
	nums := make([]int, len(split))
	for i, char := range split {
		num, _ := strconv.ParseInt(char, 10, 0)
		nums[i] = int(num)
	}
	rootNode, _ := parseNode(nums)
	fmt.Println("Part one", sumMeta(rootNode))
	fmt.Println("Part two", sumRoot(rootNode))
}
