package game

import "github.com/Manuel9550/Life/pkg/tile"

type Board struct {
	squareRows int
	squareColumns int

	tiles [][]tile.Tile
	tilesUpdate [][]tile.Tile
}

func (b *Board) initialize(width int, height int) {
	// Must fill the screen with 20x20 squares
	b.squareRows = height / 20
	b.squareColumns = width /20

	b.tiles = make([][]tile.Tile,b.squareColumns)
	b.tilesUpdate = make([][]tile.Tile,b.squareColumns)

	for x := range b.tiles {
		b.tiles[x] = make([]tile.Tile,b.squareRows)
		b.tilesUpdate[x] = make([]tile.Tile,b.squareRows)
		for y := range b.tiles[x] {
			b.tiles[x][y] = tile.Tile{Alive:false}
			b.tilesUpdate[x][y] = tile.Tile{Alive:false}
		}
	}

	b.tiles[10][10].Alive = true
	b.tiles[11][10].Alive = true
	b.tiles[12][10].Alive = true

	b.tiles[2][2].Alive = true
	b.tiles[3][3].Alive = true
	b.tiles[4][3].Alive = true
	b.tiles[4][2].Alive = true
	b.tiles[3][4].Alive = true


}

func (b *Board)  UpdateTiles() {


	for x := 0; x < b.squareColumns; x++ {
		for y := 0; y < b.squareRows; y++ {
			if b.tiles[x][y].Alive {
				// Live cells with exactly two or three live neighbours lives on to the next generation

				liveCount := b.liveCount(x,y)

				if liveCount != 2 && liveCount != 3 {
					b.tilesUpdate[x][y].Alive = false
				} else {
					b.tilesUpdate[x][y].Alive = true
				}
			} else {
				// dead cells with three live neighbours becomes a live cell
				liveCount := b.liveCount(x,y)

				if liveCount == 3 {
					b.tilesUpdate[x][y].Alive = true
				} else {
					b.tilesUpdate[x][y].Alive = false
				}
			}

		}
	}

	// Once we have the new state, copy the updated state into the state that will be rendered on screen
	for x := 0; x < b.squareColumns; x++ {
		copy(b.tiles[x],b.tilesUpdate[x])
	}
}

func (b *Board)  liveCount(x int, y int) int {
	liveCount := 0
	for i := x - 1; i <= x + 1; i++ {
		for t := y - 1; t <= y + 1; t++ {

			// We don't include the actual cell, just the neighbours!
			if i != x || t != y {
				// Make sure not to fetch cells that are out of bounds
				if i >= 0 && t >= 0 && i < b.squareColumns && t < b.squareRows {
					if b.tiles[i][t].Alive {
						liveCount += 1
					}
				}
			}
		}
	}

	return liveCount
}

