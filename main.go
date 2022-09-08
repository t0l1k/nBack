package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := ebiten.RunGame(app.NewGame()); err != nil {
		log.Fatal(err)
	}
}
