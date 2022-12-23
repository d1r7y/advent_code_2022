/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day21

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day21Cmd represents the day21 command
var day21Cmd = &cobra.Command{
	Use:   "day21",
	Short: "Day 21 of Advent of Code 2022",
	Long:  `Monkey Math`,
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
		err = day21(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day21Cmd)

	day21Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
