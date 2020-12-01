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

	// Store images for the mouse/cells
	images map[string]*ebiten.Image

	board Board // Keeps track of and updates the cells
	panel Panel // Stores the buttons the player can click on

	keyPressed   map[ebiten.Key]bool
	mousePressed bool
	paused       bool

	// Variables for dealing with the interval state
	interval     int
	NextInterval int
	timeChange   int
	maxInterval  int
	minInterval  int
	ticker       *time.Ticker

	// The current/last position of the mouse
	xPos int
	yPos int

	lastXPos int
	lastYPos int

	heightOffset int
}

func (g *Game) Init(gameWidth int, gameHeight int) error {

	g.images = make(map[string]*ebiten.Image)

	var err error

	g.images["ALIVE"], _, err = ebitenutil.NewImageFromFile("../assets/tile-green.png")
	if err != nil {
		return err
	}

	g.images["DEAD"], _, err = ebitenutil.NewImageFromFile("../assets/tile-white.png")
	if err != nil {
		return err
	}

	g.images["HIGHLIGHTED"], _, err = ebitenutil.NewImageFromFile("../assets/tile-blue.png")
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

// The main update function
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

// Checks for keyboard input, called every Update
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

// Draws the board to the screen
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

	// Draw the speed we are currently playing at

	msg := "Current Speed:"
	text.Draw(screen, msg, g.panel.font, g.panel.messageX, g.panel.messageY, color.White)

	if g.paused {
		msg = "Paused"
		text.Draw(screen, msg, g.panel.font, g.panel.pauseX, g.panel.fontY, color.White)
	} else {
		msg = strconv.Itoa(g.NextInterval) + " milliseconds"
		text.Draw(screen, msg, g.panel.font, g.panel.fontX, g.panel.fontY, color.White)
	}

}

// O
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

// Increase/Decrease the interval at which a new state is calculated
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

// Check what the current state of the cursor is
func (g *Game) checkCursor() {

	x, y := ebiten.CursorPosition()
	g.xPos, g.yPos = g.board.checkSquare(x, y)

}

// Checks for mouse input. Called on every Update
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
