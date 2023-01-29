/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/michaelbui99/chip8-go/internal/provider"
)

// hexDumpCmd represents the hexDump command
var hexDumpCmd = &cobra.Command{
	Use:   "hexDump",
	Short: "Hex dump Chip 8 RAM after loading ROM",
	Args:  cobra.ExactArgs(1),
	Long:  `Loads a ROM into the Chip 8's RAM and them hex dumps the RAM`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chip8 := provider.GetChip8Instance()
		romPath = args[0]
		fmt.Printf("LOADING: %v\n", romPath)

		if err := chip8.Load(romPath); err != nil {
			return err
		}

		fmt.Printf("%v", hex.Dump(chip8.Ram[:]))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hexDumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hexDumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hexDumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
