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

	board Board

	ticker* time.Ticker

	keyPressed map[ebiten.Key]bool // Map of key presses, to determine if the user has released a key

	paused bool

	interval time.Duration
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

	g.paused = false

	g.width = gameWidth
	g.height = gameHeight

	// Initialize the board with 20x20 tiles
	g.board.initialize(gameWidth,gameHeight)

	g.interval = 100

	// Set the timer to the standard 1 update per second
	g.ticker = time.NewTicker(g.interval * time.Millisecond)

	// Initialize the keys we are tracking
	g.keyPressed = make(map[ebiten.Key]bool)
}

func (g *Game) Update() error {

	g.checkKeys()

	// Check if enough time has passed to update the game
	for {
		select {
		case _ = <-g.ticker.C:
			// The ticker has sent a value: Perform an update
			g.board.UpdateTiles()

		default:
			return nil
		}
	}

	//return nil
}

func (g * Game) checkKeys() {

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.keyPressed[ebiten.KeySpace] = true
	} else {
		// Was this key previously pressed?
		if g.keyPressed[ebiten.KeySpace] {
			g.keyPressed[ebiten.KeySpace] = false
			g.pauseButton()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.keyPressed[ebiten.KeyLeft] = true
	} else {
		// Was this key previously pressed?
		if g.keyPressed[ebiten.KeyLeft] {
			g.keyPressed[ebiten.KeyLeft] = false
			g.updateTime(-100)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.keyPressed[ebiten.KeyRight] = true
	} else {
		// Was this key previously pressed?
		if g.keyPressed[ebiten.KeyRight] {
			g.keyPressed[ebiten.KeyRight] = false
			g.updateTime(100)
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
	g.interval += duration

	if g.interval >= 2000 {
		g.interval = 2000
	}

	if g.interval <= 100 {
		g.interval = 100
	}

	g.ticker = time.NewTicker(g.interval * time.Millisecond)
}





