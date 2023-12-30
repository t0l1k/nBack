package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/scenes/scene_select"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(scene_select.NewSceneSelectGame())
	eui.Quit()
}
