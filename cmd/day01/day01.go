/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day01

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day01Cmd represents the day01 command
var day01Cmd = &cobra.Command{
	Use:   "day01",
	Short: "Day 1 of Advent of Code 2022",
	Long:  `Calorie Counting`,
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
		err = day01(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day01Cmd)

	day01Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
