/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day05

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day05Cmd represents the day05 command
var day05Cmd = &cobra.Command{
	Use:   "day05",
	Short: "Day 5 of Advent of Code 2022",
	Long:  `Supply Stacks`,
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
		err = day05(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day05Cmd)

	day05Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
