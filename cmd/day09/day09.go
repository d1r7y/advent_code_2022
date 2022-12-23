/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day09

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day09Cmd represents the day09 command
var day09Cmd = &cobra.Command{
	Use:   "day09",
	Short: "Day 9 of Advent of Code 2022",
	Long:  `Rope Bridge`,
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
		err = day09(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day09Cmd)

	day09Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
