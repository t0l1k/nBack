package main

import (
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/ui"
)

func main() {
	ui.Init(app.NewGame())
	ui.Run(app.NewSceneToday())
	ui.Quit()
}
