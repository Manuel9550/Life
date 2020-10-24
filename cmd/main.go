package main

import (
	"github.com/Manuel9550/Life/pkg/game"
	"github.com/hajimehoshi/ebiten"
	_ "github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/color"

	"log"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Geometry Matrix")
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}


