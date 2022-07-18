package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/game"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
