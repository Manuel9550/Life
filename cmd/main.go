package main

import (
	"github.com/Manuel9550/Life/pkg/game"
	"github.com/hajimehoshi/ebiten"
	_ "github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/color"

	"log"
)

func main() {

	width := 640
	height := 480

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Conway's Game of Life")

	gameToRun := game.Game{}
	gameToRun.Init(width,height)

	if err := ebiten.RunGame(&gameToRun); err != nil {
		log.Fatal(err)
	}
}


