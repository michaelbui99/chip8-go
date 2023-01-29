package emulator

import "github.com/veandco/go-sdl2/sdl"

var sdlKeyToNativeChip8KeyMappings map[string]byte

func GetSdlKeyToChip8KeyMappings() *map[string]byte {
	if sdlKeyToNativeChip8KeyMappings == nil {
		sdlKeyToNativeChip8KeyMappings = map[string]byte{
			string(sdl.K_1): 0x1,
			string(sdl.K_2): 0x2,
			string(sdl.K_3): 0x3,
			string(sdl.K_4): 0xC,
			string(sdl.K_q): 0x4,
			string(sdl.K_w): 0x5,
			string(sdl.K_e): 0x6,
			string(sdl.K_r): 0xD,
			string(sdl.K_a): 0x7,
			string(sdl.K_s): 0x8,
			string(sdl.K_d): 0x9,
			string(sdl.K_f): 0xE,
			string(sdl.K_z): 0xA,
			string(sdl.K_x): 0x0,
			string(sdl.K_c): 0xB,
			string(sdl.K_v): 0xF,
		}
	}
	return &sdlKeyToNativeChip8KeyMappings
}

func GetChip8Key(sdlKey string) byte {
	return sdlKeyToNativeChip8KeyMappings[sdlKey]
}
