package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	const input = ``
	var lines = strings.Split(input, "\n")
	sort.Strings(lines)
	var current_guard = ""
	var fell_asleep_at = 0
	var guard_sleeping = make(map[string]int)
	var guard_minutes_asleep = make(map[string]map[int]int)
	for _, line := range lines {
		var minute, _ = strconv.ParseInt(strings.Split(strings.Split(line, ":")[1], "]")[0], 10, 0)
		var action = strings.Split(strings.Trim(strings.Split(line, "]")[1], " "), " ")[0]
		if action == "Guard" {
			current_guard = strings.Split(strings.Trim(strings.Split(line, "]")[1], " "), " ")[1]
		} else if action == "falls" {
			fell_asleep_at = int(minute)
		} else if action == "wakes" {
			guard_sleeping[current_guard] += int(minute) - fell_asleep_at
			if guard_minutes_asleep[current_guard] == nil {
				guard_minutes_asleep[current_guard] = make(map[int]int)
			}
			for i := fell_asleep_at; i < int(minute); i++ {
				guard_minutes_asleep[current_guard][i] += 1
			}
		}
	}
	var max_guard = ""
	var max_guard_number int64 = 0
	var max_sleep = 0
	for id, minutes := range guard_sleeping {
		if minutes > max_sleep {
			max_sleep = minutes
			max_guard = id
			max_guard_number, _ = strconv.ParseInt(id[1:], 10, 0)
		}
	}
	var max_minute = 0
	var max_times_asleep = 0
	for minute, times_asleep := range guard_minutes_asleep[max_guard] {
		if times_asleep > max_times_asleep {
			max_times_asleep = times_asleep
			max_minute = minute
		}
	}
	fmt.Println("Part 1:", max_guard, "*", max_minute, ":", int(max_guard_number)*max_minute)
	max_times_asleep = 0
	for id, _ := range guard_sleeping {
		for minute, times_asleep := range guard_minutes_asleep[id] {
			if times_asleep > max_times_asleep {
				max_times_asleep = times_asleep
				max_minute = minute
				max_guard_number, _ = strconv.ParseInt(id[1:], 10, 0)
			}
		}
	}
	fmt.Println("Part 2:", max_guard_number, "*", max_minute, ":", int(max_guard_number)*max_minute)

}
