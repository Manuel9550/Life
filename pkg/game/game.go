package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/png"
	"log"
	"time"
)

//var width = 640
//var height = 480
//var alive *ebiten.Image
//var dead *ebiten.Image

//var squareRows int
//var squareColumns int

//var tiles [][]tile.Tile

//var ticker* time.Ticker

type Game struct {
	width int
	height int
	alive *ebiten.Image
	dead *ebiten.Image
	highlighted *ebiten.Image

	board Board

	ticker* time.Ticker

	keyPressed map[ebiten.Key]bool // Map of key presses, to determine if the user has released a key
	mousePressed bool 
	paused bool

	interval time.Duration
	NextInterval time.Duration

	xPos int
	yPos int

	lastXPos int
	lastYPos int
	
	

}

func (g *Game)  Init(gameWidth int, gameHeight int) {

	var err error
	g.dead, _, err = ebitenutil.NewImageFromFile("assets/tile-white.png")
	if err != nil {
		log.Fatal(err)
	}

	g.alive, _, err = ebitenutil.NewImageFromFile("assets/tile-green.png")
	if err != nil {
		log.Fatal(err)
	}

	g.highlighted, _, err = ebitenutil.NewImageFromFile("assets/tile-blue.png")
	if err != nil {
		log.Fatal(err)
	}

	g.paused = false

	g.width = gameWidth
	g.height = gameHeight

	// Initialize the board with 20x20 tiles
	g.board.initialize(gameWidth,gameHeight)

	g.interval = 500
	g.NextInterval = g.interval

	// Set the timer to the standard 1 update per second
	g.ticker = time.NewTicker(g.interval * time.Millisecond)

	// Initialize the keys we are tracking
	g.keyPressed = map[ebiten.Key]bool{
		ebiten.KeySpace:false,
		ebiten.KeyLeft:false,
		ebiten.KeyRight:false,
	}

	g.xPos = 0
	g.yPos = 0

	g.lastXPos = 0
	g.lastYPos = 0
}

func (g *Game) Update() error {

	g.checkKeys()
	g.checkMouse()
	g.checkCursor()

	// Check if enough time has passed to update the game
	for {
		select {
		case _ = <-g.ticker.C:
			// The ticker has sent a value: Perform an update
			g.board.UpdateTiles()

			// Check if we should update the timer as well
			if g.NextInterval != g.interval {
				g.interval = g.NextInterval
				g.ticker = time.NewTicker(g.interval * time.Millisecond)
			}

		default:
			return nil
		}
	}

	//return nil
}

func (g * Game) checkKeys() {

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
					g.updateTime(-100)
				case ebiten.KeyRight:
					g.updateTime(100)
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
			op.GeoM.Translate(float64(x) * 20, float64(y) * 20 )

			if g.board.tiles[x][y].Alive {
				screen.DrawImage(g.alive, op)
			} else {
				screen.DrawImage(g.dead, op)
			}

		}
	}

	// whichever square the cursor is currently on, highlight it
	op.GeoM.Reset()
	op.GeoM.Translate(float64(g.xPos) * float64(g.board.squareSize), float64(g.yPos) * float64(g.board.squareSize) )
	screen.DrawImage(g.highlighted, op)


	//op.GeoM.Translate(50, 50)
	//op.GeoM.Scale(1, 1)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}


func (g* Game) pauseButton() {
	// If the timer is on, stop it. If the timer isn't on, restart it
	if g.paused {
		g.ticker = time.NewTicker(g.interval * time.Millisecond)
		g.paused = false
	} else {
		g.ticker.Stop()
		g.paused = true
	}

}

func (g* Game) updateTime(duration time.Duration) {
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
	g.xPos, g.yPos = g.board.checkSquare(x,y)

}

func (g *Game) checkMouse() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.mousePressed = true

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
			g.board.tiles[g.xPos][g.yPos].Click()
		}
	}
}





