package day16

import (
	"errors"
	"fmt"
	"strings"
)

const StartingValve = "AA"

type ValveMap map[string]*Valve

type VisitedValve map[string]bool

type Labyrinth struct {
	ValveCount                     int
	Map                            ValveMap
	Time                           int
	OpenValves                     []*Valve
	CurrentPressureReleasedPerTime int
	TotalPressureReleased          int
}

func NewLabyrinth() *Labyrinth {
	return &Labyrinth{Time: 1, Map: make(ValveMap), OpenValves: make([]*Valve, 0)}
}

func (l *Labyrinth) AddValve(v *Valve) {
	l.Map[v.Name] = v
	l.ValveCount++
}

func (l *Labyrinth) DeleteValve(valveName string) {
	delete(l.Map, valveName)
	l.ValveCount--
}

func (l *Labyrinth) FindValve(name string) *Valve {
	return l.Map[name]
}

func (l *Labyrinth) SimplifyDepthFirst(valveName string, visited VisitedValve) {
	valve := l.FindValve(valveName)

	visited[valveName] = true

	if valve.PressureRelief == 0 && valve.Name != StartingValve {
		// Here's a useless valve that isn't the starting valve.
		// Find all references to it and remove it.
		for _, st := range valve.Tunnels {
			siblingValve := l.FindValve(st.ValveName)
			newSiblingTunnels := make([]Tunnel, 0)
			for _, ct := range siblingValve.Tunnels {
				// Add existing tunnels, skipping the valve we want to remove.
				if ct.ValveName != valveName {
					newSiblingTunnels = append(newSiblingTunnels, ct)
				}
			}

			// Now add tunnels to valve's other siblings, skipping current sibling.  Tunnels need to have
			// a higher cost.
			for _, st2 := range valve.Tunnels {
				if st2.ValveName != st.ValveName {
					nt := NewTunnel(st2.ValveName, st.Cost+st2.Cost)
					newSiblingTunnels = append(newSiblingTunnels, nt)
				}
			}

			// Save updated tunnels
			siblingValve.Tunnels = newSiblingTunnels
		}

		l.DeleteValve(valveName)
	}

	for _, t := range valve.Tunnels {
		if !visited[t.ValveName] {
			l.SimplifyDepthFirst(t.ValveName, visited)
		}
	}
}

func (l *Labyrinth) Simplify() {
	visited := make(VisitedValve)

	l.SimplifyDepthFirst(StartingValve, visited)
}

func (l *Labyrinth) Tick() {
	// Choose what to do.

	// Factor travel time.  Going down paths to single nodes: account for return trip!
	l.Time++
}

type Tunnel struct {
	ValveName string
	Cost      int
}

func NewTunnel(name string, cost int) Tunnel {
	return Tunnel{ValveName: name, Cost: cost}
}

type Valve struct {
	Name           string
	PressureRelief int
	Opened         bool
	Tunnels        []Tunnel
}

func NewValve(name string, pressureRelief int) *Valve {
	return &Valve{Name: name, PressureRelief: pressureRelief, Tunnels: make([]Tunnel, 0)}
}

func (v *Valve) Open() {
	v.Opened = true
}

func ParseValveDefinition(line string) (*Valve, error) {
	// Remove the semicolon from the line.
	noSemicolon := strings.ReplaceAll(line, ";", "")

	// Turn the equal sign into space.
	noEqual := strings.ReplaceAll(noSemicolon, "=", " ")

	var valveName string
	var valvePressure int
	var tunnelValve string

	// Check for single tunnel
	count, err := fmt.Sscanf(noEqual, "Valve %s has flow rate %d tunnel leads to valve %s", &valveName, &valvePressure, &tunnelValve)
	if err == nil && count == 3 {
		newValve := NewValve(valveName, valvePressure)
		c := NewTunnel(tunnelValve, 1)
		newValve.Tunnels = append(newValve.Tunnels, c)

		return newValve, nil
	}

	// Now try multiple tunnels
	count, err = fmt.Sscanf(noEqual, "Valve %s has flow rate %d tunnels lead to valves ", &valveName, &valvePressure)
	if err != nil {
		return nil, err
	}
	if count != 2 {
		return nil, errors.New("invalid valve definition")
	}

	// Now scan the connecting valves.
	splitStr := strings.SplitAfter(noEqual, "valves ")
	if len(splitStr) != 2 {
		return nil, errors.New("invalid connecting valves")
	}

	// Remove the commas
	noCommas := strings.ReplaceAll(splitStr[1], ",", "")

	newValve := NewValve(valveName, valvePressure)
	for _, cv := range strings.Fields(noCommas) {
		c := NewTunnel(cv, 1)
		newValve.Tunnels = append(newValve.Tunnels, c)
	}

	return newValve, nil
}

func day16(fileContents string) error {
	l := NewLabyrinth()

	fmt.Println(20*3 + 33*4 + 54*8 + 76*4 + 79*3 + 81*6)

	fmt.Println(21*4 + 41*2 + 44*4 + 66*6 + 68*2 + 81*9)
	// minute 1:
	// move to open JJ
	// minute 2:
	// still moving to JJ
	// minute 3:
	// open JJ
	// minute 4:
	// JJ open 21
	// move to AA
	// minute 5:
	// JJ open 21
	// still moving to AA
	// minute 6:
	// JJ open 21
	// move to DD
	// minute 7:
	// JJ open 21
	// open DD
	// minute 8:
	// JJ/DD open 41
	// move to EE
	// minute 9:
	// JJ/DD open 41
	// open EE
	// minute 10:
	// JJ/DD/EE open 44
	// move to HH
	// minute 11:
	// JJ/DD/EE open 44
	// still moving to HH
	// minute 12:
	// JJ/DD/EE open 44
	// still moving to HH
	// minute 13:
	// JJ/DD/EE open 44
	// open HH
	// minute 14:
	// JJ/DD/EE/HH open 66
	// move to EE
	// minute 15:
	// JJ/DD/EE/HH open 66
	// still moving to EE
	// minute 16:
	// JJ/DD/EE/HH open 66
	// still moving to EE
	// minute 17:
	// JJ/DD/EE/HH open 66
	// move to DD
	// minute 18:
	// JJ/DD/EE/HH open 66
	// move to CC
	// minute 19:
	// JJ/DD/EE/HH open 66
	// open CC
	// minute 20:
	// JJ/DD/EE/HH/CC open 68
	// move to BB
	// minute 21:
	// JJ/DD/EE/HH/CC open 68
	// open BB
	// minute 22:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 23:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 24:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 25:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 26:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 27:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 28:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 29:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay
	// minute 30:
	// JJ/DD/EE/HH/CC/BB open 81
	// stay

	fmt.Println(20*2 + 23*4 + 45*6 + 47*2 + 60*4 + 81*10)
	// minute 1:
	// move to DD
	// minute 2:
	// open DD
	// minute 3:
	// DD release 20 pressure
	// move to EE
	// minute 4:
	// DD release 20 pressure
	// open EE
	// minute 5:
	// DD/EE release 23 pressure
	// move to HH
	// minute 6:
	// DD/EE release 23 pressure
	// still moving to HH
	// minute 7:
	// DD/EE release 23 pressure
	// still moving to HH
	// minute 8:
	// DD/EE release 23 pressure
	// open HH
	// minute 9:
	// DD/EE/HH release 45 pressure
	// move to EE
	// minute 10:
	// DD/EE/HH release 45 pressure
	// still moving to EE
	// minute 11:
	// DD/EE/HH release 45 pressure
	// still moving to EE
	// minute 12:
	// DD/EE/HH release 45 pressure
	// move to DD
	// minute 13:
	// DD/EE/HH release 45 pressure
	// move to CC
	// minute 14:
	// DD/EE/HH release 45 pressure
	// open CC
	// minute 15:
	// DD/EE/HH/CC release 47 pressure
	// move to BB
	// minute 16:
	// DD/EE/HH/CC release 47 pressure
	// open BB
	// minute 17:
	// DD/EE/HH/CC/BB release 60 pressure
	// move to AA
	// minute 18:
	// DD/EE/HH/CC/BB release 60 pressure
	// move to JJ
	// minute 19:
	// DD/EE/HH/CC/BB release 60 pressure
	// still move to JJ
	// minute 20:
	// DD/EE/HH/CC/BB release 60 pressure
	// open JJ
	// minute 21:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 22:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 23:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 24:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 25:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 26:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 27:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 28:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 29:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit
	// minute 30:
	// DD/EE/HH/CC/BB/JJ release 81 pressure
	// sit

	for _, line := range strings.Split(fileContents, "\n") {
		valve, err := ParseValveDefinition(line)
		if err != nil {
			return err
		}

		l.AddValve(valve)
	}

	// Simplify the graph by removing valves that don't reduce pressure.
	l.Simplify()

	fmt.Println("Valve count after simplifying ", l.ValveCount)
	return nil
}
