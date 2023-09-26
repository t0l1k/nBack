package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/ui/app"
	"github.com/t0l1k/nBack/ui/scene/today"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(today.NewSceneToday())
	eui.Quit()
}
