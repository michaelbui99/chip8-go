package emulator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/michaelbui99/chip8-go/internal/bitutil"
	"github.com/michaelbui99/chip8-go/internal/display"
	"github.com/veandco/go-sdl2/sdl"
)

const ramSize = 4096
const registerCount = 16
const memStart_dec = 512   // 0x0200
const memStart_hex = 0x200 // 0x0200

var fontSet = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F}
}

type Chip8 struct {
	Stack         *Stack
	Ram           [ramSize]byte       // 4 kB RAM
	V             [registerCount]byte // 16 general purpose registers (V0 - VF). VF is also used as a flag register
	I             uint16              // Index register
	PC            uint16              // Program Counter. First Chip 8 interpreter was also located in RAM (000 - 1FF). Start PC at 0x200 to be compatible with old programs
	DelayTimer    byte                // Should be decremented by one 60 times per second
	SoundTimer    byte                // Should be decremented by one 60 times per second
	CurrentOpCode uint16              // Current instruction to execute
	Display       *display.Display
	RomLoaded     bool
	RefreshRateMs uint32
}

func NewChip8() *Chip8 {
	chip8 := &Chip8{
		Stack:         NewStack(16),
		Ram:           [ramSize]byte{},
		V:             [registerCount]byte{},
		I:             0,
		PC:            memStart_hex,
		DelayTimer:    0x0,
		CurrentOpCode: 0x0,
		Display:       display.NewDisplay(10),
		RomLoaded:     false,
		RefreshRateMs: 60,
	}

	copy(chip8.Ram[0:len(fontSet)], fontSet)

	return chip8
}

func (c *Chip8) Start() error {
	if !c.RomLoaded {
		return errors.New("no ROM has been loaded")
	}
	defer c.Display.TearDown()

	finished := false
	for !finished {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				finished = true
			}
		}

		err := c.EmulateCycle()
		c.Display.Draw()
		if err != nil {
			return err
		}

		c.delay(1000 / c.RefreshRateMs)
	}

	return nil
}

func (c *Chip8) EmulateCycle() error {
	c.fetchNextOp()
	err := c.decodeAndExecuteOp()
	if err != nil {
		return err
	}

	return nil
}

// Loads a ROM into the Chip8's RAM. The ROM must not exceed 4kB (4096 bytes)
func (c *Chip8) Load(file string) error {
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	if fStat.Size() > ramSize {
		return errors.New("ROM exceeds Chip 8's memory")
	}

	program, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	for idx, data := range program {
		c.Ram[idx+memStart_dec] = data
	}

	c.RomLoaded = true
	return nil
}

func (c *Chip8) SetDisplaySizeFactor(sizeFactor int32) {
	c.Display = display.NewDisplay(sizeFactor)
}

// Fetches the next instruction and increments program counter by 2.
// Each instruction is 2 bytes
func (c *Chip8) fetchNextOp() {
	c.CurrentOpCode = bitutil.CombineBytes(c.Ram[c.PC], c.Ram[c.PC+1])
	c.PC += 2
}

// Executes the instruction that was loaded into the Chip8's CurrentOpCode
func (c *Chip8) decodeAndExecuteOp() error {
	// TODO: Implement the rest of the instruction set. The current ones are implemented since they are used in the IBM Logo program
	x := c.CurrentOpCode & 0x0f00
	y := c.CurrentOpCode & 0x00f0
	n := c.CurrentOpCode & 0x000f
	nn := c.CurrentOpCode & 0x00ff
	nnn := c.CurrentOpCode & 0x0fff

	//  Masking first 4 bits, since they describe the type of instruction.
	//  The rest of the opcode (x, y ,n ,nn, nnn) are just arguments
	opcodeMasked := c.CurrentOpCode & 0xf000
	switch opcodeMasked {
	case 0x0000:
		switch c.CurrentOpCode {
		case 0x00E0:
			// 0x00E0 - Clear screen
			c.Display.Clear()
			return nil

		case 0x00EE:
			// 0x00EE return from subroutine
			retAddr, err := c.Stack.Pop()
			if err != nil {
				return err
			}
			c.PC = retAddr
		}
		return nil

	case 0x1000:
		// 0x1NNN - Jump to address NNN
		c.PC = nnn
		return nil

	case 0x2000:
		// 0x2NNN - Call subroutine at memory location NNN
		c.Stack.Push(c.PC) // Push current address so we can return after subroutine
		c.PC = nnn

	case 0x3000:
		// 0x3[X][NN] - Skip one instruction (2 bytes) if value of VX == NN
		vx := c.Ram[x>>8]
		if vx == byte(nn) {
			c.PC += 2
		}
		return nil

	case 0x4000:
		// 0x4[X][NN] - Skip one instruction if value of VX != NN
		vx := c.Ram[x>>8]
		if vx != byte(nn) {
			c.PC += 2
		}
		return nil

	case 0x5000:
		// 0x5XY0 - Skip one instruction if value of VX == value of VY
		vx := c.Ram[x>>8]
		vy := c.Ram[y>>4]
		if vx == vy {
			c.PC += 2
		}
		return nil

	case 0x9000:
		// 0x9XY0 - Skip one instruction if value of VX != value of VY
		vx := c.Ram[x>>8]
		vy := c.Ram[y>>4]
		if vx != vy {
			c.PC += 2
		}
		return nil

	case 0x6000:
		// 0x6[X][NN] - Sets register VX to NN
		c.Ram[x>>8] = byte(nn)
		return nil

	case 0x7000:
		// 0x7[X][NN] - Add value NN to register VX
		c.Ram[x>>8] = c.Ram[x>>8] + byte(nn)
		return nil

	case 0x8000:
		// 0x8XY* - Logical arithemetic instructions. The specific operations i determined by the last 4 bits
		vx := c.Ram[x>>8]
		vy := c.Ram[y>>4]
		switch c.CurrentOpCode & 0x000F {
		case 0x0000:
			// 0x8XY0 - Set VX to value of VY
			c.Ram[x>>8] = vy
			return nil

		case 0x0001:
			// 0x8XY1 - set VX to VX OR VY
			c.Ram[x>>8] = vx | vy
			return nil

		case 0x0002:
			// 0x8XY2 - Set VX to VX AND VY
			c.Ram[x>>8] = vx & vy
			return nil

		case 0x0003:
			// 0x8XY3 - Set VX to VX XOR VY
			c.Ram[x>>8] = (vx | vy) & ^(vx & vy)
			return nil

		case 0x0004:
			// 0x8XY4 - Set VX to VX ADD VY. Set VF to 1 if carry else 0
			res := vx + vy

			if res > 0xFF {
				c.Ram[0xF] = 1
			} else {
				c.Ram[0xF] = 0
			}

			c.Ram[x>>8] = res
			return nil

		case 0x0005:
			// 0x8XY5 - Set VX to VX SUBTRACT VY. Set VF to 0 if underflow else 1
			if vx >= vy {
				c.Ram[0xF] = 1
			} else {
				c.Ram[0xF] = 0
			}

			c.Ram[x>>8] = vx - vy
			return nil

		case 0x0007:
			// 0x8XY5 - Set VX to VY SUBTRACT VX. Set VF to 0 if underflow else 1
			if vy >= vx {
				c.Ram[0xF] = 1
			} else {
				c.Ram[0xF] = 0
			}

			c.Ram[x>>8] = vy - vx
			return nil

		case 0x0006:
			// 0x8XY6 - Right shift VX one bit. Set VF to 1 if bit shifted out was 1 else 0
			vxBit0 := vx & 0x1
			if vxBit0 == 1 {
				c.Ram[0xF] = 1
			} else {
				c.Ram[0xF] = 0
			}

			c.Ram[x>>8] = vx >> 1
			return nil

		case 0x000E:
			// 0x8XYE - Left shift VX one bit. Set VF to 1 if bit shifted out was 1 else 0
			vxBit7 := vx & 0b10000000
			if vxBit7 == 1 {
				c.Ram[0xF] = 1
			} else {
				c.Ram[0xF] = 0
			}

			c.Ram[x>>8] = vx << 1
			return nil
		}
		return nil

	case 0xA000:
		// 0xANNN - Set index register I to address NNN
		c.I = nnn
		return nil

	case 0xB000:
		// 0xBNNN - Jump with offset. Jump to address NNN + V0
		c.PC = nnn + uint16(c.Ram[0x0])
		return nil

	case 0xC000:
		// 0xCXNN - Generate random number, rn, between 0 and NN, do rn AND NN and store result in VX
		rn := rand.Intn(int(nn) + 1)
		res := byte(rn) & c.Ram[x>>8]
		c.Ram[x>>8] = res
		return nil

	case 0xD000:
		// 0xD[X][Y][N]  - Draw N pixels tall sprite at (VX, VY) from the memory location that I is holding to the screen
		d_x := c.Ram[x>>8] % byte(display.Width)
		d_y := c.Ram[y>>4] % byte(display.Height)
		c.Ram[0xf] = 0x0

		for i := uint16(0); i < n; i++ {
			pixel := c.Ram[c.I+i]
			for j := uint16(0); j < 8; j++ {
				if (pixel & (0x80 >> j)) != 0 {
					if c.Display.Instance[(d_y + byte(i))][d_x+byte(j)] == 1 {
						c.Ram[0xf] = 1
					}

					c.Display.Instance[(d_y + byte(i))][d_x+byte(j)] ^= 1
				}
			}
		}
		c.Display.Draw()
		return nil

	case 0xF000:
		switch c.CurrentOpCode & 0x00FF {
		case 0x0007:
			// 0xFX07 - Set VX to current value of delay timer
			c.Ram[x>>8] = c.DelayTimer
			return nil

		case 0x0015:
			// 0xFX15 - Set Delay Timer to value in VX
			c.DelayTimer = c.Ram[x>>8]
			return nil

		case 0x0018:
			// 0xFX18 - Set Sound Timer to value in VX
			c.SoundTimer = c.Ram[x>>8]
			return nil
		case 0x001E:
			// 0xFX1E - Set I to I + VX
			c.I = c.I + uint16(c.Ram[x>>8])
			return nil
		}
		return nil

	default:
		return fmt.Errorf("UKNOWN INSTRUCTION: %X PC: %v\n%v", opcodeMasked, c.PC, hex.Dump(c.Ram[:]))
	}

	return nil
}

func (c *Chip8) delay(ms uint32) {
	sdl.Delay(ms)
}
