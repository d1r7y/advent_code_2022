package day15

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Sensor struct {
	Position       Point
	Beacon         Point
	BeaconDistance int
}

type SensorList []*Sensor

func (s SensorList) Len() int {
	return len(s)
}

func (s SensorList) Less(i, j int) bool {
	if s[i].Position.X < s[j].Position.X {
		return true
	}
	if s[i].Position.X > s[j].Position.X {
		return false
	}

	// X values are equal.

	if s[i].Position.Y < s[j].Position.Y {
		return true
	}

	return false
}

func (s SensorList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Network struct {
	Min     Point
	Max     Point
	Sensors SensorList
}

func (n *Network) ClosestSensor(p Point) *Sensor {
	closestDistance := math.MaxInt
	var closestSensor *Sensor

	for _, s := range n.Sensors {
		d := ManhattanDistance(s.Position, p)

		if d < closestDistance {
			closestDistance = d
			closestSensor = s
		}
	}

	return closestSensor
}

func (n *Network) SensorIntersection(p Point) SensorList {
	sensors := make(SensorList, 0)

	for _, s := range n.Sensors {
		d := ManhattanDistance(s.Position, p)

		if d <= s.BeaconDistance {
			sensors = append(sensors, s)
		}
	}

	return sensors
}

func (n *Network) InvalidBeaconLocations(row int) []Point {
	// Since sensors can overlap, duplicate points can appear.  Create a map
	// with the points as keys to keep them unique.
	locationMap := make(map[Point]bool)

	for i := n.Min.X; i <= n.Max.X; i++ {
		p := Point{i, row}

		sensors := n.SensorIntersection(p)
		for _, s := range sensors {
			if p != s.Beacon {
				d := ManhattanDistance(p, s.Position)
				if d <= s.BeaconDistance {
					locationMap[p] = true
				}
			}
		}
	}

	locations := make([]Point, 0)

	// Extract the keys.
	for p := range locationMap {
		locations = append(locations, p)
	}

	return locations
}

func (n *Network) PossibleBeaconLocations() []Point {
	locations := make([]Point, 0)

	for i := n.Min.Y; i <= n.Max.Y; i++ {
	NextLocation:
		for j := n.Min.X; j <= n.Max.X; j++ {
			p := Point{j, i}

			for _, s := range n.Sensors {
				d := ManhattanDistance(s.Position, p)

				if d <= s.BeaconDistance {
					j = s.Position.X + s.BeaconDistance - AbsoluteDifference(i, s.Position.Y)
					continue NextLocation
				}
			}

			locations = append(locations, p)
		}
	}

	return locations
}

func AbsoluteDifference(a, b int) int {
	if a < b {
		return b - a
	}

	return a - b
}

func ManhattanDistance(p1, p2 Point) int {
	return AbsoluteDifference(p1.X, p2.X) + AbsoluteDifference(p1.Y, p2.Y)
}

func NewSensor(position Point, beacon Point) *Sensor {
	s := &Sensor{Position: position, Beacon: beacon}
	s.BeaconDistance = ManhattanDistance(s.Position, s.Beacon)
	return s
}

func ParseSensorLine(line string) *Sensor {
	var sensorX, sensorY int
	var beaconX, beaconY int

	count, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensorX, &sensorY, &beaconX, &beaconY)
	if err != nil {
		log.Panic(err)
	}

	if count != 4 {
		log.Panic("invalid sensor line")
	}

	return NewSensor(Point{sensorX, sensorY}, Point{beaconX, beaconY})
}

func ParseSensors(fileContents string) SensorList {
	sensors := make(SensorList, 0)

	for _, line := range strings.Split(fileContents, "\n") {
		sensor := ParseSensorLine(line)
		sensors = append(sensors, sensor)
	}

	return sensors
}

func ParseNetwork(fileContents string, tightBounds bool) *Network {
	sensors := ParseSensors(fileContents)
	n := &Network{Sensors: sensors}

	n.Min = Point{math.MaxInt, math.MaxInt}
	n.Max = Point{math.MinInt, math.MinInt}

	for _, sensor := range sensors {
		if tightBounds {
			if sensor.Position.X < n.Min.X {
				n.Min.X = sensor.Position.X
			}
			if sensor.Position.Y < n.Min.Y {
				n.Min.Y = sensor.Position.Y
			}

			if sensor.Position.X > n.Max.X {
				n.Max.X = sensor.Position.X
			}
			if sensor.Position.Y > n.Max.Y {
				n.Max.Y = sensor.Position.Y
			}
		} else {
			if sensor.Position.X-sensor.BeaconDistance < n.Min.X {
				n.Min.X = sensor.Position.X - sensor.BeaconDistance
			}
			if sensor.Position.Y-sensor.BeaconDistance < n.Min.Y {
				n.Min.Y = sensor.Position.Y - sensor.BeaconDistance
			}

			if sensor.Position.X+sensor.BeaconDistance > n.Max.X {
				n.Max.X = sensor.Position.X + sensor.BeaconDistance
			}
			if sensor.Position.Y+sensor.BeaconDistance > n.Max.Y {
				n.Max.Y = sensor.Position.Y + sensor.BeaconDistance
			}
		}
	}

	return n
}

func day15(fileContents string) error {
	// Part 1: Given a sensor report containing sensor locations and the closest beacons to
	// them, which locations, in a given row, cannot contain a beacon?
	n := ParseNetwork(fileContents, false)

	row := 2000000
	invalidLocations := n.InvalidBeaconLocations(row)

	fmt.Printf("Beacons cannot be in %d locations in row %d.\n", len(invalidLocations), row)

	// Part 2: Given a sensor report containing sensor locations and the closest beacons to
	// them, there is only a single location where the distress beacon can be.  You can calculate
	// its tuning frequency by multiplying its x coordinate by 4000000 and adding its y coordinate.
	n = ParseNetwork(fileContents, true)

	validLocations := n.PossibleBeaconLocations()

	if len(validLocations) > 1 {
		fmt.Printf("Unexpected number of possible beacon locations: %d\n", len(validLocations))
		return nil
	}

	fmt.Printf("Distress beacon tuning frequency %d.\n", validLocations[0].X*4000000+validLocations[0].Y)

	return nil
}
