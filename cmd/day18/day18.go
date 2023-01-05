/*
Copyright © 2023 Cameron Esfahani <dirty@mac.com>
*/
package day18

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day18Cmd represents the day18 command
var day18Cmd = &cobra.Command{
	Use:   "day18",
	Short: "Day 18 of Advent of Code 2022",
	Long:  `Boiling Boulders`,
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
		err = day18(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day18Cmd)

	day18Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
