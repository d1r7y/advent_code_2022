/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day10

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day10Cmd represents the day10 command
var day10Cmd = &cobra.Command{
	Use:   "day10",
	Short: "Day 10 of Advent of Code 2022",
	Long:  `Cathode-Ray Tube`,
	Run: func(cmd *cobra.Command, args []string) {
		df, err := os.Open(inputPath)
		if err != nil {
			log.Fatal(err)
		}

		defer df.Close()

		fileContent, err := io.ReadAll(df)
		if err != nil {
			log.Fatal(err)
		}
		err = day10(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day10Cmd)

	day10Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
