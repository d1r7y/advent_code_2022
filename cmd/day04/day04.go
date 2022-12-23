/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day04

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day04Cmd represents the day04 command
var day04Cmd = &cobra.Command{
	Use:   "day04",
	Short: "Day 4 of Advent of Code 2022",
	Long:  `Camp Cleanup`,
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
		err = day04(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day04Cmd)

	day04Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
