package game

import (
	"github.com/Manuel9550/Life/pkg/tile"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/png"
	"log"
)

var width = 640
var height = 480
var alive *ebiten.Image
var dead *ebiten.Image

var tiles [][]tile.Tile

type Game struct{}

func init() {
	var err error
	dead, _, err = ebitenutil.NewImageFromFile("assets/tile-white.png")
	if err != nil {
		log.Fatal(err)
	}

	alive, _, err = ebitenutil.NewImageFromFile("assets/tile-blue.png")
	if err != nil {
		log.Fatal(err)
	}



	// Must fill the screen with 20x20 squares
	squareRows := width / 20
	squareColumns := height /20

	tiles = make([][]tile.Tile,squareRows)

	for x := range tiles {
		tiles[x] = make([]tile.Tile,squareColumns)
		for y := range tiles[x] {
			tiles[x][y] = tile.Tile{Alive:false}
		}
	}

	// Initialize the tiles to be all empty
	for x := 0; x < squareRows; x++ {
		for y := 0; y < squareColumns; y++ {
			tiles[x][y] = tile.Tile{false}
		}
	}

	// Test the array!
	tiles[10][10].Alive = true
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Must fill the screen with 20x20 squares
	squareRows := width / 20
	squareColumns := height /20

	op := &ebiten.DrawImageOptions{}

	for x := 0; x < squareRows; x++ {
		for y := 0; y < squareColumns; y++ {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(x) * 20, float64(y) * 20 )

			if tiles[x][y].Alive {
				screen.DrawImage(alive, op)
			} else {
				screen.DrawImage(dead, op)
			}

		}
	}


	//op.GeoM.Translate(50, 50)
	//op.GeoM.Scale(1, 1)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
