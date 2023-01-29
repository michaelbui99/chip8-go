package provider

import "github.com/michaelbui99/chip8-go/internal/emulator"

var instance *emulator.Chip8

func GetChip8Instance() *emulator.Chip8 {
	if instance == nil {
		instance = emulator.NewChip8()
	}

	return instance
}
