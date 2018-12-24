package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type AttackType string

type Group struct {
	count       int
	hp          int
	initiative  int
	weaknesses  []AttackType
	immunities  []AttackType
	attackPower int
	attackType  AttackType
	isReindeer  bool
	name        string
}

func (g *Group) effectivePower() int {
	return g.count * g.attackPower
}

func (g *Group) isImmune(t AttackType) bool {
	for _, attackType := range g.immunities {
		if attackType == t {
			return true
		}
	}
	return false
}

func (g *Group) isWeak(t AttackType) bool {
	for _, attackType := range g.weaknesses {
		if attackType == t {
			return true
		}
	}
	return false
}

func (g *Group) attackDamage(otherGroup *Group) int {
	if otherGroup == nil || otherGroup.isImmune(g.attackType) {
		return 0
	}
	damage := g.effectivePower()
	if otherGroup.isWeak(g.attackType) {
		return damage * 2
	}
	return damage
}

func (g *Group) attack(otherGroup *Group) {
	damage := g.attackDamage(otherGroup) / otherGroup.hp
	if otherGroup.count < damage {
		damage = otherGroup.count
	}
	otherGroup.count -= damage
	// fmt.Println(g.name, "attacks group", otherGroup.name, "for", damage, "killing", damage)
}

type Powers []*Group

func (s Powers) Len() int {
	return len(s)
}
func (s Powers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Powers) Less(i, j int) bool {
	l1 := s[i]
	l2 := s[j]
	if l1.effectivePower() == l2.effectivePower() {
		return l1.initiative > l2.initiative
	}
	return l1.effectivePower() > l2.effectivePower()
}

type Initiatives []*Group

func (s Initiatives) Len() int {
	return len(s)
}
func (s Initiatives) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Initiatives) Less(i, j int) bool {
	l1 := s[i]
	l2 := s[j]
	return l1.initiative > l2.initiative
}

func getGroups(input []byte, boost int) []*Group {
	lines := strings.Split(string(input), "\n")
	units := []*Group{}
	isReindeer := true
	counter := 1
	for _, line := range lines {
		if len(line) < 10 {
			continue
		}
		if strings.Contains(line, ":") {
			if line == "Infection:" {
				isReindeer = false
				counter = 1
			}
			continue
		}
		parts := strings.Split(line, " ")
		count, _ := strconv.Atoi(parts[0])
		hp, _ := strconv.Atoi(parts[4])
		initiative, _ := strconv.Atoi(parts[len(parts)-1])
		attackPower, _ := strconv.Atoi(parts[len(parts)-6])
		attackType := AttackType(parts[len(parts)-5])
		weaknesses := []AttackType{}
		immunities := []AttackType{}
		// Get weaknesses and immunities
		if strings.Contains(line, "(") {
			typesPart := strings.Split(strings.Split(line, ")")[0], "(")[1]
			for _, part := range strings.Split(typesPart, ";") {
				part = strings.Trim(part, " ")
				types := strings.Split(strings.Split(part, "to ")[1], ", ")
				kind := strings.Split(part, " ")[0]
				for _, t := range types {
					if kind == "weak" {
						weaknesses = append(weaknesses, AttackType(t))
					} else {
						immunities = append(immunities, AttackType(t))
					}
				}
			}
		}
		var name string
		if isReindeer {
			name = fmt.Sprintf("Immune System Group %d", counter)
		} else {
			name = fmt.Sprintf("Infection Group %d", counter)
		}
		if isReindeer {
			attackPower += boost
		}
		units = append(units, &Group{count, hp, initiative, weaknesses, immunities, attackPower, attackType, isReindeer, name})
		counter++
	}
	return units
}

func fight(units []*Group) (bool, int) {
	prev := 0
	for {
		// target selection
		sort.Sort(Powers(units))
		isTargeted := make(map[*Group]bool)
		selectedTarget := make(map[*Group]*Group)
		for _, thisgroup := range units {
			var willAttack *Group
			var willAttackDamage int
			for _, othergroup := range units {
				if thisgroup.isReindeer != othergroup.isReindeer && !isTargeted[othergroup] && othergroup.count > 0 {
					attackDamage := thisgroup.attackDamage(othergroup)
					if attackDamage == 0 {
						continue
					}
					// fmt.Println(thisgroup.name, "would deal", othergroup.name, attackDamage)
					if attackDamage >= willAttackDamage {
						if attackDamage == willAttackDamage {
							if othergroup.effectivePower() >= willAttack.effectivePower() {
								if othergroup.effectivePower() == willAttack.effectivePower() {
									if othergroup.initiative > willAttack.initiative {
										willAttack = othergroup
										willAttackDamage = attackDamage
									}
								} else {
									willAttack = othergroup
									willAttackDamage = attackDamage
								}
							}
						} else {
							willAttack = othergroup
							willAttackDamage = attackDamage
						}
					}
				}
			}
			if willAttack != nil {
				// fmt.Println(thisgroup.name, "(", thisgroup.effectivePower(), ") will attack", willAttack.name)
				isTargeted[willAttack] = true
				selectedTarget[thisgroup] = willAttack
			}
		}
		// attacking phase
		sort.Sort(Initiatives(units))
		attacksMade := false
		for _, group := range units {
			if group.count == 0 || selectedTarget[group] == nil {
				continue
			}
			attacksMade = true
			group.attack(selectedTarget[group])
		}
		if !attacksMade {
			break
		}
		totalAlive := 0
		for _, group := range units {
			totalAlive += group.count
		}
		if prev == totalAlive {
			break
		}
		prev = totalAlive
	}
	didImmuneSystemWin := true
	totalAlive := 0
	for _, group := range units {
		totalAlive += group.count
		if group.count > 0 {
			didImmuneSystemWin = didImmuneSystemWin && group.isReindeer
		}
	}
	return didImmuneSystemWin, totalAlive
}

func main() {
	buf, _ := ioutil.ReadFile("day-24.txt")
	units := getGroups(buf, 0)
	_, winners := fight(units)
	fmt.Println("Part One:", winners)
	boost := 1
	for {
		immuneSystemDidWin, n := fight(getGroups(buf, boost))
		if immuneSystemDidWin {
			fmt.Println("Part Two:", n)
			break
		}
		boost++
	}
}
