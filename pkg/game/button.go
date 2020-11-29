package game

import (
	"github.com/hajimehoshi/ebiten"
)

type Button struct {
	x int
	y int

	height int
	width  int

	image *ebiten.Image

	frameScaleWidth  float64
	frameScaleHeight float64

	scaleX float64
	scaleY float64
}

func (b *Button) IsPressed(x int, y int) bool {
	if x >= b.x && x <= (b.x+int(b.width)) && y >= b.y && y <= (b.y+int(b.height)) {
		return true
	} else {
		return false
	}
}

func (b *Button) SetOp(frameOp *ebiten.DrawImageOptions, buttonOp *ebiten.DrawImageOptions, frameImage *ebiten.Image) {

	// Draw the button frame before the button image
	frameOp.GeoM.Reset()

	frameOp.GeoM.Scale(b.frameScaleWidth, b.frameScaleHeight)
	frameOp.GeoM.Translate(float64(b.x), float64(b.y))

	//

	buttonOp.GeoM.Reset()
	buttonOp.GeoM.Scale(b.scaleX, b.scaleY)
	buttonOp.GeoM.Translate(float64(b.x)+(float64(b.width)/4), float64(b.y)+(float64(b.height)/4))
}
