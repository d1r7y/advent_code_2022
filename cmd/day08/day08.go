/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day08

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day08Cmd represents the day08 command
var day08Cmd = &cobra.Command{
	Use:   "day08",
	Short: "Day 8 of Advent of Code 2022",
	Long:  `Treetop Tree House`,
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
		err = day08(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day08Cmd)

	day08Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
