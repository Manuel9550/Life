package game

import (
	"github.com/hajimehoshi/ebiten"
)

type Button struct {
	x int
	y int

	height float64
	width  float64

	image *ebiten.Image

	frameScaleWidth  float64
	frameScaleHeight float64
}

func (b *Button) Init(locationX int, locationY int, buttonWidth float64, buttonHeight float64, buttonType string, buttonImage *ebiten.Image) {
	b.x = locationX
	b.y = locationY

	b.width = buttonWidth
	b.height = buttonHeight

	b.image = buttonImage

}

func (b *Button) IsPressed(x int, y int) bool {
	if x >= b.x && x <= (b.x+int(b.width)) && y >= b.y && y <= (b.y+int(b.height)) {
		return true
	} else {
		return false
	}
}
