package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Panel struct {
	images map[string]*ebiten.Image

	buttons map[string]*Button

	font  font.Face
	fontX int
	fontY int

	messageX int
	messageY int
	pauseX   int
}

type buttonParam struct {
	buttonWidth  int
	buttonHeight int
}

type panelInput struct {
	gameHeight   int
	gameWidth    int
	heightOffset int
	button       buttonParam
}

func (p *Panel) initialize(panelStats *panelInput) error {

	p.images = make(map[string]*ebiten.Image)

	var err error
	p.images["DEAD"], _, err = ebitenutil.NewImageFromFile("assets/tile-white.png")
	if err != nil {
		return err
	}

	p.images["ALIVE"], _, err = ebitenutil.NewImageFromFile("assets/tile-green.png")
	if err != nil {
		return err
	}

	p.images["HIGHLIGHTED"], _, err = ebitenutil.NewImageFromFile("assets/tile-blue.png")
	if err != nil {
		return err
	}

	p.images["PLAY"], _, err = ebitenutil.NewImageFromFile("assets/play-icon.png")
	if err != nil {
		return err
	}

	p.images["PAUSE"], _, err = ebitenutil.NewImageFromFile("assets/pause-icon.png")
	if err != nil {
		return err
	}

	p.images["SLOWER"], _, err = ebitenutil.NewImageFromFile("assets/slower-icon.png")
	if err != nil {
		return err
	}

	p.images["FASTER"], _, err = ebitenutil.NewImageFromFile("assets/faster-icon.png")
	if err != nil {
		return err
	}

	p.images["FRAME"], _, err = ebitenutil.NewImageFromFile("assets/button-frame.png")
	if err != nil {
		return err
	}

	bw, bh := p.images["PLAY"].Size()
	fw, fh := p.images["FRAME"].Size()

	p.buttons = make(map[string]*Button)

	buttonYLocation := panelStats.gameHeight - int(panelStats.button.buttonHeight) - panelStats.heightOffset

	// Scale the buttons and frames so they are the proper size on the screen
	heightScale := (float64(panelStats.button.buttonHeight) + float64(panelStats.heightOffset)) / float64(bh) * 0.5
	widthScale := float64(panelStats.button.buttonWidth) / float64(bw) * 0.5

	// Each button will have it's own frame behind it
	frameScaleHeight := (float64(panelStats.button.buttonHeight) + float64(panelStats.heightOffset)) / float64(fh)
	frameScaleWidth := float64(panelStats.button.buttonWidth) / float64(fw)

	p.buttons["PLAY"] = &Button{
		x:                0,
		y:                buttonYLocation,
		height:           panelStats.button.buttonHeight,
		width:            panelStats.button.buttonWidth,
		image:            p.images["PLAY"],
		frameScaleWidth:  frameScaleWidth,
		frameScaleHeight: frameScaleHeight,
		scaleX:           widthScale,
		scaleY:           heightScale,
	}

	p.buttons["SLOWER"] = &Button{
		x:                panelStats.button.buttonWidth,
		y:                buttonYLocation,
		height:           panelStats.button.buttonHeight,
		width:            panelStats.button.buttonWidth,
		image:            p.images["SLOWER"],
		frameScaleWidth:  frameScaleWidth,
		frameScaleHeight: frameScaleHeight,
		scaleX:           widthScale,
		scaleY:           heightScale,
	}

	p.buttons["FASTER"] = &Button{
		x:                panelStats.button.buttonWidth * 4,
		y:                buttonYLocation,
		height:           panelStats.button.buttonHeight,
		width:            panelStats.button.buttonWidth,
		image:            p.images["FASTER"],
		frameScaleWidth:  frameScaleWidth,
		frameScaleHeight: frameScaleHeight,
		scaleX:           widthScale,
		scaleY:           heightScale,
	}

	// Font
	const dpi = 72
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	p.font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	if err != nil {
		return err
	}

	p.fontX = int(float64(panelStats.button.buttonWidth) * 2.25) // start halfway between the end of the their button area
	p.fontY = int(buttonYLocation) + int(float64(panelStats.button.buttonHeight)/1.25)

	p.messageX = int(float64(panelStats.button.buttonWidth) * 2.3)
	p.messageY = int(buttonYLocation) + int(panelStats.button.buttonHeight/2)

	p.pauseX = int(float64(panelStats.button.buttonWidth) * 2.70)

	return nil
}
