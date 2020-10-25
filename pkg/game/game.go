package game

import (
	"github.com/Manuel9550/Life/pkg/tile"
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

	squareRows int
	squareColumns int

	tiles [][]tile.Tile
	tilesUpdate [][]tile.Tile

	ticker* time.Ticker
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

	g.width = gameWidth
	g.height = gameHeight

	// Must fill the screen with 20x20 squares
	g.squareRows = g.height / 20
	g.squareColumns = g.width /20

	g.tiles = make([][]tile.Tile,g.squareColumns)
	g.tilesUpdate = make([][]tile.Tile,g.squareColumns)

	for x := range g.tiles {
		g.tiles[x] = make([]tile.Tile,g.squareRows)
		g.tilesUpdate[x] = make([]tile.Tile,g.squareRows)
		for y := range g.tiles[x] {
			g.tiles[x][y] = tile.Tile{Alive:false}
			g.tilesUpdate[x][y] = tile.Tile{Alive:false}
		}
	}


	// Test the array!
	g.tiles[10][10].Alive = true
	g.tiles[11][10].Alive = true
	g.tiles[12][10].Alive = true

	// Set the timer to the standard 1 update per second
	g.ticker = time.NewTicker(1000 * time.Millisecond)
}

func (g *Game) Update() error {


	// Check if enough time has passed to update the game
	for {
		select {
		case _ = <-g.ticker.C:
			// The ticker has sent a value: Perform an update
			g.UpdateTiles()

		default:
			return nil
		}
	}

	//return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Must fill the screen with 20x20 squares

	op := &ebiten.DrawImageOptions{}

	for x := 0; x < g.squareColumns; x++ {
		for y := 0; y < g.squareRows; y++ {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(x) * 20, float64(y) * 20 )

			if g.tiles[x][y].Alive {
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

func (g *Game)  UpdateTiles() {


	for x := 0; x < g.squareColumns; x++ {
		for y := 0; y < g.squareRows; y++ {
			if g.tiles[x][y].Alive {
				// Live cells with exactly two or three live neighbours lives on to the next generation

				liveCount := g.liveCount(x,y)

				if liveCount != 2 && liveCount != 3 {
					g.tilesUpdate[x][y].Alive = false
				} else {
					g.tilesUpdate[x][y].Alive = true
				}
			} else {
				// dead cells with three live neighbours becomes a live cell
				liveCount := g.liveCount(x,y)

				if liveCount == 3 {
					g.tilesUpdate[x][y].Alive = true
				} else {
					g.tilesUpdate[x][y].Alive = false
				}
			}

		}
	}

	// Once we have the new state, copy the updated state into the state that will be rendered on screen
	for x := 0; x < g.squareColumns; x++ {
		copy(g.tiles[x],g.tilesUpdate[x])
	}
}

func (g *Game)  liveCount(x int, y int) int {
	liveCount := 0
	for i := x - 1; i <= x + 1; i++ {
		for t := y - 1; t <= y + 1; t++ {

			// We don't include the actual cell, just the neighbours!
			if i != x || t != y {
				// Make sure not to fetch cells that are out of bounds
				if i >= 0 && t >= 0 && i < g.squareColumns && t < g.squareRows {
					if g.tiles[i][t].Alive {
						liveCount += 1
					}
				}
			}
		}
	}

	return liveCount
}

