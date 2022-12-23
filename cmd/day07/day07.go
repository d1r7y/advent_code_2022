/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day07

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day07Cmd represents the day07 command
var day07Cmd = &cobra.Command{
	Use:   "day07",
	Short: "Day 7 of Advent of Code 2022",
	Long:  `No Space Left On Device`,
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
		err = day07(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day07Cmd)

	day07Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
