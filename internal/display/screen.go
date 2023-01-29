package display

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Dimensions
const (
	Height int32 = 32
	Width  int32 = 64
)

type Display struct {
	Instance   [Height][Width]byte
	Window     *sdl.Window
	Renderer   *sdl.Renderer
	SizeFactor int32
}

func NewDisplay(sizeFactor int32) *Display {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("chip8-go", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, Width*sizeFactor, Height*sizeFactor, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}

	// Set background to black
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()

	return &Display{
		Instance:   [Height][Width]byte{},
		Window:     window,
		Renderer:   renderer,
		SizeFactor: sizeFactor,
	}
}

func (d *Display) Clear() {
	for i := range d.Instance {
		for j := range d.Instance[i] {
			d.Instance[i][j] = 0x0
		}
	}
}

func (d *Display) Draw(screenState [Height][Width]byte) {
	// Clear window
	d.Renderer.SetDrawColor(0, 0, 0, 0)
	d.Renderer.Clear()

	// Draw screen
	for i, row := range d.Instance {
		for j := range row {
			pixelOn := screenState[i][j] != 0
			if pixelOn {
				// Draw with white if pixel on
				d.Renderer.SetDrawColor(255, 255, 255, 255)
			} else {
				// Draw black if pixel not on
				d.Renderer.SetDrawColor(0, 0, 0, 255)
			}

			d.Renderer.FillRect(&sdl.Rect{
				X: int32(j) * d.SizeFactor,
				Y: int32(i) * d.SizeFactor,
				W: d.SizeFactor,
				H: d.SizeFactor,
			})
		}
	}

	d.render()
}

func (d *Display) render() {
	d.Renderer.Present()
}

func (d *Display) TearDown() {
	if d.Window != nil {
		d.Window.Destroy()
	}

	if d.Renderer != nil {
		d.Renderer.Destroy()
	}
}
