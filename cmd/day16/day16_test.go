package day16

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValveDefinition(t *testing.T) {
	type testCase struct {
		str           string
		expectedErr   bool
		expectedValve Valve
	}

	testCases := []testCase{
		{"Valve AA has flow rate=0; tunnels lead to valves DD, II, BB", false, Valve{"AA", 0, false, []Tunnel{{"DD", 1}, {"II", 1}, {"BB", 1}}}},
		{"Valve BB has flow rate=13; tunnels lead to valves CC, AA", false, Valve{"BB", 13, false, []Tunnel{{"CC", 1}, {"AA", 1}}}},
		{"Valve HH has flow rate=22; tunnel leads to valve GG", false, Valve{"HH", 22, false, []Tunnel{{"GG", 1}}}},
	}
	for _, test := range testCases {
		valve, err := ParseValveDefinition(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			assert.Equal(t, test.expectedValve, *valve)
		}
	}
}

func TestSimplify(t *testing.T) {
	valveDefinitions := `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`

	l := NewLabyrinth()

	for _, line := range strings.Split(valveDefinitions, "\n") {
		valve, err := ParseValveDefinition(line)
		assert.NoError(t, err)

		l.AddValve(valve)
	}

	// Simplify the graph by removing valves that don't reduce pressure.
	l.Simplify()
}
