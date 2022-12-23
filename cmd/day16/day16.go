/*
Copyright © 2022 Cameron Esfahani <dirty@mac.com>
*/
package day16

import (
	"io"
	"log"
	"os"

	"github.com/d1r7y/advent_2022/cmd"
	"github.com/spf13/cobra"
)

// day16Cmd represents the day16 command
var day16Cmd = &cobra.Command{
	Use:   "day16",
	Short: "Day 16 of Advent of Code 2022",
	Long:  `Proboscidea Volcanium`,
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
		err = day16(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var inputPath string

func init() {
	cmd.RootCmd.AddCommand(day16Cmd)

	day16Cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file")
}
