package day06

import (
	"fmt"
)

type Datastream string

func (d Datastream) GetPacketMarkerStart() int {
	const StartOfPacketMarkerLength = 4

	characters := make([]byte, 0)

NextCharacter:
	for i, c := range d {
		if len(characters) == StartOfPacketMarkerLength {
			// Remove the first character
			characters = characters[1:]
		}

		characters = append(characters, byte(c))

		// Do we have enough characters to look for the packet marker?
		if len(characters) == StartOfPacketMarkerLength {
			uniqueMap := make(map[byte]bool)

			for _, c := range characters {
				if _, ok := uniqueMap[c]; ok {
					continue NextCharacter
				}

				uniqueMap[c] = true
			}

			return i + 1
		}
	}

	return -1
}

func (d Datastream) GetMessageMarkerStart() int {
	const StartOfMessageMarkerLength = 14

	characters := make([]byte, 0)

NextCharacter:
	for i, c := range d {
		if len(characters) == StartOfMessageMarkerLength {
			// Remove the first character
			characters = characters[1:]
		}

		characters = append(characters, byte(c))

		// Do we have enough characters to look for the message marker?
		if len(characters) == StartOfMessageMarkerLength {
			uniqueMap := make(map[byte]bool)

			for _, c := range characters {
				if _, ok := uniqueMap[c]; ok {
					continue NextCharacter
				}

				uniqueMap[c] = true
			}

			return i + 1
		}
	}

	return -1
}

func day06(fileContents string) error {
	// Part 1: What is the offset of the first valid packet marker?  Packets need 4 unique characters.
	ds := Datastream(string(fileContents))
	if validOffset := ds.GetPacketMarkerStart(); validOffset > 0 {
		fmt.Printf("Valid packet marker starting at offset %d\n", validOffset)
	} else {
		fmt.Print("No valid packet marker found")
	}

	// Part 2: What is the offset of the first valid message marker?  Messages need 14 unique characters.
	if validOffset := ds.GetMessageMarkerStart(); validOffset > 0 {
		fmt.Printf("Valid message marker starting at offset %d\n", validOffset)
	} else {
		fmt.Print("No valid message marker found")
	}

	return nil
}
