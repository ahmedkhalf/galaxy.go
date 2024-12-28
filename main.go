//go:build js && wasm

package main

import (
	"image/color"
	"log"
	"math"
	"strconv"
	"syscall/js"

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
	scale        int
	universe     Universe
	nextPoints   chan [][2]int
}

func (g *ScreenSaver) Update() error {
	scale := float64(g.screenWidth+g.screenHeight) / 8.0
	midx, midy := g.screenWidth/2, g.screenHeight/2
	g.universe.UpdateGalaxy(scale, midx, midy)
	return nil
}

func (g *ScreenSaver) Draw(screen *ebiten.Image) {
	points := <-g.nextPoints
	for _, point := range points {
		color := hueToRGB(float64(g.universe.Galaxies[0].Galcol) / (COLORBASE - 1))
		screen.Set(point[0], point[1], color)
	}
}

func (g *ScreenSaver) generateNextPoints() {
	for {
		var points [][2]int
		for _, galaxy := range g.universe.Galaxies {
			for _, point := range galaxy.Newpoints {
				points = append(points, point)
			}
		}
		g.nextPoints <- points
	}
}

func (_ *ScreenSaver) Layout(w, h int) (int, int) {
	panic("ebitengine version must support LayoutF()")
}

func (g *ScreenSaver) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor() / float64(g.scale)

	// Ebiten uses ceil internally to calculate the screen size, based on v2.8.5
	g.screenWidth = int(math.Ceil(logicWinWidth * scale))
	g.screenHeight = int(math.Ceil(logicWinHeight * scale))

	return logicWinWidth * scale, logicWinHeight * scale
}

func GetURLParameter(name string) string {
	window := js.Global().Get("window")
	searchParams := window.Get("location").Get("search")
	urlSearchParams := js.Global().Get("URLSearchParams").New(searchParams)
	return urlSearchParams.Call("get", name).String()
}

func getScale() int {
	scale := GetURLParameter("scale")
	i, err := strconv.Atoi(scale)
	if err != nil {
		return 1
	}
	return i
}

func getFPS() int {
	fps := GetURLParameter("fps")
	i, err := strconv.Atoi(fps)
	if err != nil || i < 1 || i > 60 {
		return 30
	}
	return i
}

func main() {
	initialWidth := 800
	initialHeight := 600
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Resizable Black Window with Center Red Pixel")
	ebiten.SetTPS(getFPS())

	game := &ScreenSaver{
		screenWidth:  initialWidth,
		screenHeight: initialHeight,
		universe:     InitGalaxy(),
		scale:        getScale(),
		nextPoints:   make(chan [][2]int, 1),
	}
	go game.generateNextPoints()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
