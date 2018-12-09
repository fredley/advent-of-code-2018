package main

import (
	"fmt"
	"sort"
	"strings"
)

type Edge struct {
	from string
	to   string
}

type Worker struct {
	time      int
	workingOn string
}

func predecessors(char string, edges []Edge) []string {
	var preds = make([]string, 0)
	for _, edge := range edges {
		if edge.to == char {
			preds = append(preds, edge.from)
		}
	}
	return preds
}

func successors(char string, edges []Edge) []string {
	var succs = make([]string, 0)
	for _, edge := range edges {
		if edge.from == char {
			succs = append(succs, edge.to)
		}
	}
	return succs
}

func findRoots(edges []Edge) []string {
	var roots = make([]string, 0)
	for _, edge := range edges {
		if len(predecessors(edge.from, edges)) == 0 {
			var unique = true
			for _, root := range roots {
				if root == edge.from {
					unique = false
					break
				}
			}
			if unique {
				roots = append(roots, edge.from)
			}
		}
	}
	return roots
}

func allPredecessorsIn(char string, seen string, edges []Edge) bool {
	preds := predecessors(char, edges)
	for _, pred := range preds {
		if !strings.Contains(seen, pred) {
			return false
		}
	}
	return true
}

func filterSeenUnique(ls []string, seen string) []string {
	var filtered = make([]string, 0)
	var seenHere = ""
	for _, char := range ls {
		if !strings.Contains(seen, char) && !strings.Contains(seenHere, char) {
			filtered = append(filtered, char)
			seenHere += char
		}
	}
	return filtered
}

func filterReady(ls []string, done string, doing string, edges []Edge) []string {
	var filtered = make([]string, 0)
	for _, char := range ls {
		if allPredecessorsIn(char, done, edges) && !strings.Contains(done+doing, char) {
			filtered = append(filtered, char)
		}
	}
	return filtered
}

func extra_time(c string) int {
	return int(rune(c[0])) - 64
}

func main() {
	const input = ``
	split := strings.Split(input, "\n")
	edges := make([]Edge, len(split))
	for i, line := range split {
		words := strings.Split(line, " ")
		edges[i] = Edge{words[1], words[7]}
	}
	roots := findRoots(edges)
	sort.Strings(roots)
	var current = roots[0]
	var tasks = roots[0]
	var options = roots
	var reachable = roots
	for {
		reachable = filterSeenUnique(append(reachable, successors(current, edges)...), tasks)
		options = filterReady(reachable, tasks, "", edges)
		sort.Strings(options)
		if len(options) > 0 {
			current, options = options[0], options[1:]
			tasks += current
		} else {
			break
		}
	}
	fmt.Println("Part One:", tasks)
	taskLength := 60
	var time = 0
	var workers = make([]Worker, 5)
	var done = ""
	var doing = ""
	options = roots
	reachable = roots
	for i, _ := range workers {
		workers[i].workingOn = "."
	}
	for {
		var working = false
		// finish any current jobs
		for i, worker := range workers {
			if worker.time == 0 && worker.workingOn != "." {
				done += worker.workingOn
				doing = strings.Replace(doing, worker.workingOn, "", -1)
				reachable = filterSeenUnique(append(reachable, successors(worker.workingOn, edges)...), done)
				workers[i].workingOn = "."
			}
		}
		options = filterReady(reachable, done, doing, edges)
		// assign any new tasks
		for i, worker := range workers {
			if worker.time == 0 && len(options) > 0 {
				if len(options) > 0 {
					workers[i].workingOn, options = string(options[0]), options[1:]
					doing += workers[i].workingOn
					options = filterReady(reachable, done, doing, edges)
					workers[i].time = taskLength - 1 + extra_time(workers[i].workingOn)
					working = true
				} else {
					workers[i].workingOn = "."
				}
			} else if worker.time != 0 {
				workers[i].time--
				working = true
			}
		}
		if !working {
			break
		}
		time++
	}
	fmt.Println("Part Two:", time)
}
