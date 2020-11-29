package main

import (
	_ "image/color"

	"github.com/Manuel9550/Life/pkg/game"
	"github.com/hajimehoshi/ebiten"
	_ "github.com/hajimehoshi/ebiten/ebitenutil"

	"log"
)

func main() {

	width := 640
	height := 580

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Conway's Game of Life")

	gameToRun := game.Game{}
	err := gameToRun.Init(width, height)

	if err != nil {
		log.Fatal(err)
	}

	if err = ebiten.RunGame(&gameToRun); err != nil {
		log.Fatal(err)
	}
}
