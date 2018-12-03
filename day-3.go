package main

import (
	"fmt"
	"strconv"
	"strings"
)

type cut struct {
	id        int
	left_edge int
	top_edge  int
	width     int
	height    int
}

func main() {
	const input = ``
	const SIZE = 1000
	var lines = strings.Split(input, "\n")
	var cuts = make([]cut, len(lines))
	var material = make([][]int, SIZE)
	for row := range material {
		material[row] = make([]int, SIZE)
	}
	for i, line := range lines {
		var splits = strings.Split(strings.Trim(line, " \t"), " ")
		var id = splits[0]
		var edge = splits[2]
		var dims = splits[3]
		var parsed_id, _ = strconv.ParseInt(string(id[1:]), 0, 0)
		var parsed_left_edge, _ = strconv.ParseInt(strings.Split(edge, ",")[0], 0, 0)
		var parsed_top_edge, _ = strconv.ParseInt(strings.Replace(strings.Split(edge, ",")[1], ":", "", 1), 0, 0)
		var parsed_width, _ = strconv.ParseInt(strings.Split(dims, "x")[0], 0, 0)
		var parsed_height, _ = strconv.ParseInt(strings.Split(dims, "x")[1], 0, 0)
		cuts[i] = cut{
			id:        int(parsed_id),
			left_edge: int(parsed_left_edge),
			top_edge:  int(parsed_top_edge),
			width:     int(parsed_width),
			height:    int(parsed_height),
		}
		for w := parsed_left_edge; w < parsed_left_edge+parsed_width; w++ {
			for h := parsed_top_edge; h < parsed_top_edge+parsed_height; h++ {
				material[w][h] += 1
			}
		}
	}
	var overlapped = 0
	for row := range material {
		for col := range material {
			if material[row][col] > 1 {
				overlapped += 1
			}
		}
	}
	fmt.Println("Part one:", overlapped)
	for _, cut := range cuts {
		var overlapped = false
		for w := cut.left_edge; w < cut.left_edge+cut.width; w++ {
			for h := cut.top_edge; h < cut.top_edge+cut.height; h++ {
				if material[w][h] > 1 {
					overlapped = true
					break
				}
			}
			if overlapped {
				break
			}
		}
		if !overlapped {
			fmt.Println("Part two:", cut.id)
			break
		}
	}

}
