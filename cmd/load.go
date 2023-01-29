/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/michaelbui99/chip8-go/internal/provider"
)

// loadCmd represents the load command
var romPath string

var sizeFactor int32 = 10
var loadCmd = &cobra.Command{
	Use:   "load [path to rom]",
	Short: "Load Chip-8 ROM and start Chip-8",
	Long:  `Load Chip-8 ROM and start Chip-8. The ROM must not exceed 4kB`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chip8 := provider.GetChip8Instance()
		chip8.SetDisplaySizeFactor(sizeFactor)
		romPath = args[0]
		fmt.Printf("ROM PATH: %v\n", romPath)

		err := chip8.Load(romPath)
		if err != nil {
			return err
		}

		err = chip8.Start()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	loadCmd.Flags().Int32Var(&sizeFactor, "sizefactor", 10, "Increase display size of some factor. Original size (size factor 1) of Chip-8 display is 32x64")
	rootCmd.AddCommand(loadCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
