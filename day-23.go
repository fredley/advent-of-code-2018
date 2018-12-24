package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Bot struct {
	x int
	y int
	z int
	r int
}

func (b *Bot) inRange(otherBot *Bot) bool {
	return math.Abs(float64(b.x-otherBot.x))+math.Abs(float64(b.y-otherBot.y))+math.Abs(float64(b.z-otherBot.z)) <= float64(b.r)
}

func countInRange(bots []Bot, bot *Bot) int {
	botsInRange := 0
	for _, b := range bots {
		if bot.inRange(&b) {
			botsInRange++
		}
	}
	return botsInRange
}

func countCanSee(bots []Bot, bot *Bot) int {
	botsInRange := 0
	for _, b := range bots {
		if b.inRange(bot) {
			botsInRange++
		}
	}
	return botsInRange
}

func main() {
	buf, _ := ioutil.ReadFile("day-23.txt")
	lines := strings.Split(string(buf), "\n")
	bots := make([]Bot, len(lines))
	var strongestBotIdx = 0
	var strongestR = 0
	minX := 0
	maxX := 0
	minY := 0
	maxY := 0
	minZ := 0
	maxZ := 0
	for i, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(strings.Split(parts[0], "<")[1])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(strings.Split(parts[2], ">")[0])
		r, _ := strconv.Atoi(strings.Split(parts[3], "=")[1])
		if r > strongestR {
			strongestR = r
			strongestBotIdx = i
		}
		if x > maxX {
			maxX = x
		}
		if x < minX {
			minX = x
		}
		if y > maxY {
			maxY = y
		}
		if y < minY {
			minY = y
		}
		if z > maxZ {
			maxZ = z
		}
		if z < minZ {
			minZ = z
		}
		bots[i] = Bot{x, y, z, r}
	}
	fmt.Println("Part One:", countInRange(bots, &bots[strongestBotIdx]))
	fmt.Println("Ranges x <", minX, maxX, "> y <", minY, maxY, "> z <", minZ, maxZ, ">")
	dist := 1
	for {
		if dist < maxX-minX {
			dist *= 2
		} else {
			break
		}
	}
	maxInRange := 0
	bestLocation := Bot{0, 0, 0, 0}
	best_val := 9999999999999.0
	for {
		for x := minX; x <= maxX; x += dist {
			for y := minY; y <= maxY; y += dist {
				for z := minZ; z <= maxZ; z += dist {
					testBot := Bot{x, y, z, 0}
					canSee := countCanSee(bots, &testBot)
					if canSee > maxInRange {
						bestLocation = testBot
						maxInRange = canSee
						best_val = math.Abs(float64(testBot.x)) + math.Abs(float64(testBot.y)) + math.Abs(float64(testBot.z))
					} else if canSee == maxInRange {
						if math.Abs(float64(testBot.x))+math.Abs(float64(testBot.y))+math.Abs(float64(testBot.z)) < best_val {
							best_val = math.Abs(float64(testBot.x)) + math.Abs(float64(testBot.y)) + math.Abs(float64(testBot.z))
							bestLocation = testBot
							maxInRange = canSee
						}
					}
				}
			}
		}
		if dist == 1 {
			fmt.Println("Part Two:", int(best_val))
			break
		} else {
			minX = bestLocation.x - dist
			maxX = bestLocation.x + dist
			minY = bestLocation.y - dist
			maxY = bestLocation.y + dist
			minZ = bestLocation.z - dist
			maxZ = bestLocation.z + dist
			dist /= 2
		}
	}
}
