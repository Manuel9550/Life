package tile

type Tile struct{
	Alive bool
}

func (t *Tile) Click() {
	if t.Alive {
		t.Alive = false
	} else {
		t.Alive = true
	}
}
