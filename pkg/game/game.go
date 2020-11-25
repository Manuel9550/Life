package game

import (
	"image/color"
	_ "image/png"
	"log"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	width  int
	height int
	images map[string]*ebiten.Image

	board Board

	ticker *time.Ticker

	keyPressed   map[ebiten.Key]bool // Map of key presses, to determine if the user has released a key
	mousePressed bool
	paused       bool

	interval     int
	NextInterval int
	timeChange   int

	xPos int
	yPos int

	lastXPos int
	lastYPos int

	buttonHeight float64
	buttonWidth  float64

	heightOffset int

	heightScale float64
	widthScale  float64

	buttons map[string]*Button

	font  font.Face
	fontX int
	fontY int

	messageX int
	messageY int

	pauseX int
}

func (g *Game) Init(gameWidth int, gameHeight int) {

	g.images = make(map[string]*ebiten.Image)

	var err error
	g.images["DEAD"], _, err = ebitenutil.NewImageFromFile("assets/tile-white.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["ALIVE"], _, err = ebitenutil.NewImageFromFile("assets/tile-green.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["HIGHLIGHTED"], _, err = ebitenutil.NewImageFromFile("assets/tile-blue.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["PLAY"], _, err = ebitenutil.NewImageFromFile("assets/play-icon.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["PAUSE"], _, err = ebitenutil.NewImageFromFile("assets/pause-icon.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["SLOWER"], _, err = ebitenutil.NewImageFromFile("assets/slower-icon.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["FASTER"], _, err = ebitenutil.NewImageFromFile("assets/faster-icon.png")
	if err != nil {
		log.Fatal(err)
	}

	g.images["FRAME"], _, err = ebitenutil.NewImageFromFile("assets/button-frame.png")
	if err != nil {
		log.Fatal(err)
	}

	bw, bh := g.images["PLAY"].Size()
	fw, fh := g.images["FRAME"].Size()

	g.paused = true

	g.width = gameWidth
	g.height = gameHeight

	// The game will need a bottom panel for buttons.
	g.buttonHeight = float64(gameHeight) * 0.15
	g.buttonWidth = float64(gameWidth) * 0.20

	// Initialize the board with 20x20 tiles
	g.heightOffset = g.board.initialize(gameWidth, gameHeight-int(g.buttonHeight))

	// Scale the buttons and frames so they are the proper size on the screen
	g.heightScale = (g.buttonHeight + float64(g.heightOffset)) / float64(bh) * 0.5
	g.widthScale = g.buttonWidth / float64(bw) * 0.5

	frameScaleHeight := (g.buttonHeight + float64(g.heightOffset)) / float64(fh)
	frameScaleWidth := g.buttonWidth / float64(fw)

	g.interval = 500
	g.NextInterval = g.interval

	// Set the timer to the standard 1 update per 500 milliseconds
	g.ticker = time.NewTicker(time.Duration(g.interval) * time.Millisecond)
	g.ticker.Stop()

	// How many milliseconds the interval is increased/decreased by
	g.timeChange = 100

	// Initialize the keys we are tracking
	g.keyPressed = map[ebiten.Key]bool{
		ebiten.KeySpace: false,
		ebiten.KeyLeft:  false,
		ebiten.KeyRight: false,
	}

	g.xPos = 0
	g.yPos = 0

	g.lastXPos = 0
	g.lastYPos = 0

	g.buttons = make(map[string]*Button)

	buttonHeight := g.height - int(g.buttonHeight) - g.heightOffset

	g.buttons["PLAY"] = &Button{
		x:                0,
		y:                buttonHeight,
		height:           g.buttonHeight,
		width:            g.buttonWidth,
		image:            g.images["PLAY"],
		frameScaleWidth:  frameScaleWidth,
		frameScaleHeight: frameScaleHeight,
	}

	g.buttons["SLOWER"] = &Button{
		x:                int(g.buttonWidth),
		y:                buttonHeight,
		height:           g.buttonHeight,
		width:            g.buttonWidth,
		image:            g.images["SLOWER"],
		frameScaleWidth:  frameScaleWidth,
		frameScaleHeight: frameScaleHeight,
	}

	g.buttons["FASTER"] = &Button{
		x:                int(g.buttonWidth * 4),
		y:                buttonHeight,
		height:           g.buttonHeight,
		width:            g.buttonWidth,
		image:            g.images["FASTER"],
		frameScaleWidth:  frameScaleWidth,
		frameScaleHeight: frameScaleHeight,
	}

	// Font

	const dpi = 72
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	g.font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	g.fontX = int(g.buttonWidth * 2.25) // start halfway between the end of the their button area
	g.fontY = buttonHeight + int(g.buttonHeight/1.25)

	g.messageX = int(g.buttonWidth * 2.3)
	g.messageY = buttonHeight + int(g.buttonHeight/2)

	g.pauseX = int(g.buttonWidth * 2.70)

}

func (g *Game) Update() error {

	g.checkKeys()
	g.checkMouse()
	g.checkCursor()

	// Check if enough time has passed to update the game
	//for {
	select {
	case _ = <-g.ticker.C:
		// The ticker has sent a value: Perform an update
		g.board.UpdateTiles()

		// Check if we should update the timer as well
		if g.NextInterval != g.interval {
			g.interval = g.NextInterval
			g.ticker = time.NewTicker(time.Duration(g.interval) * time.Millisecond)
		}

	default:
		return nil
	}
	//}

	return nil
}

func (g *Game) checkKeys() {

	for key, isPressed := range g.keyPressed {
		if ebiten.IsKeyPressed(key) {
			g.keyPressed[key] = true
		} else {
			// Was this key previously pressed?
			if isPressed {
				g.keyPressed[key] = false

				switch pressedKey := key; pressedKey {
				case ebiten.KeySpace:
					g.pauseButton()
				case ebiten.KeyLeft:
					g.updateTime(g.timeChange)
				case ebiten.KeyRight:
					g.updateTime(-g.timeChange)
				}

			}
		}
	}

}

func (g *Game) Draw(screen *ebiten.Image) {

	// Must fill the screen with 20x20 squares

	op := &ebiten.DrawImageOptions{}

	for x := 0; x < g.board.squareColumns; x++ {
		for y := 0; y < g.board.squareRows; y++ {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(x)*20, float64(y)*20)

			if g.board.tiles[x][y].Alive {
				screen.DrawImage(g.images["ALIVE"], op)
			} else {
				screen.DrawImage(g.images["DEAD"], op)
			}

		}
	}

	// whichever square the cursor is currently on, highlight it
	op.GeoM.Reset()
	op.GeoM.Translate(float64(g.xPos)*float64(g.board.squareSize), float64(g.yPos)*float64(g.board.squareSize))
	screen.DrawImage(g.images["HIGHLIGHTED"], op)

	// Draw the buttons

	for _, button := range g.buttons {

		// Draw the button frame before the button image
		op.GeoM.Reset()

		op.GeoM.Scale(button.frameScaleWidth, button.frameScaleHeight)
		op.GeoM.Translate(float64(button.x), float64(button.y))

		screen.DrawImage(g.images["FRAME"], op)

		op.GeoM.Reset()
		//op.GeoM.Scale(2, g.buttonHeight + float64(g.heightOffset) / float64(h))
		op.GeoM.Scale(g.widthScale, g.heightScale)
		op.GeoM.Translate(float64(button.x)+(g.buttonWidth/4), float64(button.y)+(g.buttonHeight/4))

		screen.DrawImage(button.image, op)
	}

	msg := "Current Speed:"
	text.Draw(screen, msg, g.font, g.messageX, g.messageY, color.White)

	// Draw the speed we are currently playing at
	if g.paused {
		msg = "Paused"
		text.Draw(screen, msg, g.font, g.pauseX, g.fontY, color.White)
	} else {
		msg = strconv.Itoa(g.NextInterval) + " milliseconds"
		text.Draw(screen, msg, g.font, g.fontX, g.fontY, color.White)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func (g *Game) pauseButton() {
	// If the timer is on, stop it. If the timer isn't on, restart it
	if g.paused {
		g.ticker = time.NewTicker(time.Duration(g.interval) * time.Millisecond)
		g.paused = false
		g.buttons["PLAY"].image = g.images["PAUSE"]
	} else {
		g.ticker.Stop()
		g.paused = true
		g.buttons["PLAY"].image = g.images["PLAY"]
	}

}

func (g *Game) updateTime(duration int) {
	g.NextInterval += duration

	if g.NextInterval >= 2000 {
		g.NextInterval = 2000
	}

	if g.NextInterval <= 100 {
		g.NextInterval = 100
	}
}

func (g *Game) checkCursor() {
	// compare the cursor position to any onscreen objects
	x, y := ebiten.CursorPosition()

	// check if the cursor is over a button
	g.xPos, g.yPos = g.board.checkSquare(x, y)

}

func (g *Game) checkMouse() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		if g.mousePressed == false {
			g.mousePressed = true

			g.lastXPos = g.xPos
			g.lastYPos = g.yPos
		}

		// If the user is click-dragging, then allow them to change the tile they passed into
		if g.lastXPos != g.xPos || g.lastYPos != g.yPos {

			// Only want to allow this for fulling tiles
			if !g.board.tiles[g.xPos][g.yPos].Alive {
				g.board.tiles[g.xPos][g.yPos].Click()
				g.lastXPos = g.xPos
				g.lastYPos = g.yPos
			}
		}

	} else {
		// Check if the mouse was previously pressed
		if g.mousePressed {
			g.mousePressed = false

			// Check if the user clicked any of the buttons
			x, y := ebiten.CursorPosition()
			buttonClicked := false
			for buttonName, button := range g.buttons {
				if button.IsPressed(x, y) {
					switch buttonName {
					case "PLAY":
						g.pauseButton()

					case "FASTER":
						g.updateTime(-g.timeChange)
					case "SLOWER":
						g.updateTime(g.timeChange)

					}

					buttonClicked = true
					break
				}
			}
			if !buttonClicked {
				g.board.tiles[g.xPos][g.yPos].Click()
			}

		}
	}
}
