package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	app = getApp()
	app.Push(NewSceneToday())
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
