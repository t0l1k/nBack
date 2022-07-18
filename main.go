package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
