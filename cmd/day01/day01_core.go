/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day01

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

func ParseElfCalorieList(text string) ([]int, error) {

	calories := make([]int, 0)
	currentCalories := 0

	if text == "" {
		return calories, nil
	}

	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			// New elf data.  Only add to our list if we parsed
			// any calories
			if currentCalories > 0 {
				calories = append(calories, currentCalories)
				currentCalories = 0
			}
		} else {
			var calorie int

			count, err := fmt.Sscanln(line, &calorie)
			if err != nil {
				return []int{}, err
			}

			if count != 1 {
				return calories, errors.New("unexpected line in input")
			}
			currentCalories += calorie
		}
	}

	if currentCalories > 0 {
		calories = append(calories, currentCalories)
	}

	return calories, nil
}

func day01(fileContent string) error {
	calorieList, err := ParseElfCalorieList(string(fileContent))
	if err != nil {
		log.Fatal(err)
	}

	// Sort in decreasing order
	sort.Sort(sort.Reverse(sort.IntSlice(calorieList)))

	// Part 1: What's the most calories a single elf is carrying?
	if len(calorieList) > 0 {
		log.Printf("Maximum elf calories: %d\n", calorieList[0])
	} else {
		log.Println("No elf calories in input file.")
		return nil
	}

	// Part 2: How many calories are the top three elves carrying?
	if len(calorieList) >= 3 {
		log.Printf("Maximum calories from top 3 elves: %d\n", calorieList[0]+calorieList[1]+calorieList[2])
	} else {
		log.Println("Not enough elves in input file.")
	}

	return nil
}
