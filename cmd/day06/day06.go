/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day06

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day06Cmd represents the day06 command
var day06Cmd = &cobra.Command{
	Use:   "day06",
	Short: "Day 6 of Advent of Code 2022",
	Long:  `Tuning Trouble`,
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
		err = day06(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day06Cmd)

	day06Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
