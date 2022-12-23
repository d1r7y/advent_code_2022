/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day17

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day17Cmd represents the day17 command
var day17Cmd = &cobra.Command{
	Use:   "day17",
	Short: "Day 17 of Advent of Code 2022",
	Long:  `Pyroclastic Flow`,
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
		err = day17(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day17Cmd)

	day17Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
