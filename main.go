package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func min3(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), c)
}

func hueToRGB(h float64) color.RGBA {
	kr := math.Mod(5+h*6, 6)
	kg := math.Mod(3+h*6, 6)
	kb := math.Mod(1+h*6, 6)

	r := 1 - math.Max(min3(kr, 4-kr, 1), 0)
	g := 1 - math.Max(min3(kg, 4-kg, 1), 0)
	b := 1 - math.Max(min3(kb, 4-kb, 1), 0)

	return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
}

type ScreenSaver struct {
	screenWidth  int
	screenHeight int
	universe     Universe
}

func (g *ScreenSaver) Update() error {
	scale := float64(g.screenWidth+g.screenHeight) / 8.0
	midx, midy := g.screenWidth/2, g.screenHeight/2
	g.universe.UpdateGalaxy(scale, midx, midy)
	return nil
}

func (g *ScreenSaver) Draw(screen *ebiten.Image) {
	for _, galaxy := range g.universe.Galaxies {
		for _, point := range galaxy.Newpoints {
			color := hueToRGB(float64(galaxy.Galcol) / (COLORBASE - 1))
			screen.Set(point[0], point[1], color)
		}
	}
}

func (_ *ScreenSaver) Layout(w, h int) (int, int) {
	panic("ebitengine version must support LayoutF()")
}

func (g *ScreenSaver) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()

	// Ebiten uses ceil internally to calculate the screen size, based on v2.8.5
	g.screenWidth = int(math.Ceil(logicWinWidth * scale))
	g.screenHeight = int(math.Ceil(logicWinHeight * scale))

	return logicWinWidth * scale, logicWinHeight * scale
}

func main() {
	initialWidth := 800
	initialHeight := 600
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Resizable Black Window with Center Red Pixel")
	ebiten.SetTPS(30)

	game := &ScreenSaver{screenWidth: initialWidth, screenHeight: initialHeight, universe: InitGalaxy()}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
