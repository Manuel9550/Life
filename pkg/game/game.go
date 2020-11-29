package game

import (
	"image/color"
	_ "image/png"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
)

type Game struct {
	width  int
	height int
	images map[string]*ebiten.Image

	board Board
	panel Panel

	ticker *time.Ticker

	keyPressed   map[ebiten.Key]bool // Map of key presses, to determine if the user has released a key
	mousePressed bool
	paused       bool

	interval     int
	NextInterval int
	timeChange   int
	maxInterval  int
	minInterval  int

	xPos int
	yPos int

	lastXPos int
	lastYPos int

	heightOffset int
}

func (g *Game) Init(gameWidth int, gameHeight int) error {

	g.images = make(map[string]*ebiten.Image)

	var err error
	g.images["DEAD"], _, err = ebitenutil.NewImageFromFile("assets/tile-white.png")
	if err != nil {
		return err
	}

	g.images["ALIVE"], _, err = ebitenutil.NewImageFromFile("assets/tile-green.png")
	if err != nil {
		return err
	}

	g.images["HIGHLIGHTED"], _, err = ebitenutil.NewImageFromFile("assets/tile-blue.png")
	if err != nil {
		return err
	}

	g.paused = true

	g.width = gameWidth
	g.height = gameHeight

	// The game will need a bottom panel for buttons.
	buttonHeight := int(float64(gameHeight) * 0.15)
	buttonWidth := int(float64(gameWidth) * 0.20)

	// Initialize the board with 20x20 tiles
	g.heightOffset = g.board.initialize(gameWidth, gameHeight-int(buttonHeight))

	g.interval = 500
	g.NextInterval = g.interval

	g.minInterval = 100
	g.maxInterval = 2000

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

	// Initialize the panel
	buttonStats := buttonParam{
		buttonWidth:  buttonWidth,
		buttonHeight: buttonHeight,
	}

	panelStats := panelInput{
		gameHeight:   gameHeight,
		gameWidth:    gameWidth,
		heightOffset: g.heightOffset,
		button:       buttonStats,
	}

	g.panel = Panel{}
	err = g.panel.initialize(&panelStats)

	if err != nil {
		return err
	}

	return nil

}

func (g *Game) Update() error {

	g.checkKeys()
	g.checkMouse()
	g.checkCursor()

	// Check if enough time has passed to update the game

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
					g.updateTime(-g.timeChange)
				case ebiten.KeyRight:
					g.updateTime(g.timeChange)
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
			op.GeoM.Translate(float64(x*g.board.squareSize), float64(y*g.board.squareSize))

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
	buttonOp := &ebiten.DrawImageOptions{}
	frameOp := &ebiten.DrawImageOptions{}
	for _, button := range g.panel.buttons {

		button.SetOp(frameOp, buttonOp, g.panel.images["FRAME"])
		screen.DrawImage(g.panel.images["FRAME"], frameOp)
		screen.DrawImage(button.image, buttonOp)
	}

	msg := "Current Speed:"
	text.Draw(screen, msg, g.panel.font, g.panel.messageX, g.panel.messageY, color.White)

	// Draw the speed we are currently playing at
	if g.paused {
		msg = "Paused"
		text.Draw(screen, msg, g.panel.font, g.panel.pauseX, g.panel.fontY, color.White)
	} else {
		msg = strconv.Itoa(g.NextInterval) + " milliseconds"
		text.Draw(screen, msg, g.panel.font, g.panel.fontX, g.panel.fontY, color.White)
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
		g.panel.buttons["PLAY"].image = g.panel.images["PAUSE"]
	} else {
		g.ticker.Stop()
		g.paused = true
		g.panel.buttons["PLAY"].image = g.panel.images["PLAY"]
	}

}

func (g *Game) updateTime(duration int) {

	if g.paused {
		return
	}
	g.NextInterval += duration

	if g.NextInterval >= g.maxInterval {
		g.NextInterval = g.maxInterval
	}

	if g.NextInterval <= g.minInterval {
		g.NextInterval = g.minInterval
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
			for buttonName, button := range g.panel.buttons {
				if button.IsPressed(x, y) {
					switch buttonName {
					case "PLAY":
						g.pauseButton()

					case "FASTER":
						g.updateTime(g.timeChange)
					case "SLOWER":
						g.updateTime(-g.timeChange)

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
