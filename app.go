package main

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

var app *App

func getApp() (a *App) {
	if app == nil {
		db := &Db{}
		db.Setup()
		pref := NewSettings()
		var w, h int
		if pref.fullScreen {
			w, h = ebiten.ScreenSizeInFullscreen()
		} else {
			w, h = fitWindowSize()
		}
		ebiten.SetWindowTitle("nBack")
		ebiten.SetFullscreen(pref.fullScreen)
		ebiten.SetWindowSize(w, h)
		rect := ui.NewRect([]int{0, 0, w, h})
		scns := []ui.Scene{}
		a = &App{
			fullScreen:  pref.fullScreen,
			rect:        rect,
			lastDt:      -1,
			scenes:      scns,
			db:          db,
			theme:       NewTheme(),
			preferences: pref,
		}
		log.Printf("App init: screen size:[%v, %v]", w, h)
	} else {
		a = app
	}
	return a
}

func fitWindowSize() (w int, h int) {
	ww, hh := ebiten.ScreenSizeInFullscreen()
	k := 10
	w, h = 180*k, 320*k
	for ww <= w || hh <= h {
		k -= 1
		w, h = 200*k, 320*k
	}
	return w, h
}

type App struct {
	fullScreen   bool
	rect         *ui.Rect
	scenes       []ui.Scene
	currentScene ui.Scene
	lastDt       int
	db           *Db
	theme        *Theme
	preferences  *Setting
}

func (a *App) GetScreenSize() (w, h int) {
	return a.rect.W, a.rect.H
}

func (a *App) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		a.Pop()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF11) {
		a.fullScreen = !a.fullScreen
		var w, h int
		if a.fullScreen {
			ebiten.SetFullscreen(a.fullScreen)
			w, h = ebiten.ScreenSizeInFullscreen()
		} else {
			w, h = fitWindowSize()
		}
		ebiten.SetFullscreen(a.fullScreen)
		ebiten.SetWindowSize(w, h)
		a.rect = ui.NewRect([]int{0, 0, w, h})
		for _, scene := range a.scenes {
			scene.Resize()
		}

	}
	a.currentScene.Update(a.getTick())
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(a.theme.bg)
	a.currentScene.Draw(screen)
}

func (a *App) Layout(oW, oH int) (int, int) {
	return oW, oH
}

func (a *App) Push(sc ui.Scene) {
	a.scenes = append(a.scenes, sc)
	a.currentScene = sc
	a.currentScene.Entered()
	log.Println("Scene push")
}

func (a *App) Pop() {
	if len(a.scenes) > 0 {
		a.currentScene.Quit()
		idx := len(a.scenes) - 1
		a.scenes = a.scenes[:idx]
		log.Printf("App Pop Scene Quit done.")
	}
	if len(a.scenes) > 0 {
		a.currentScene = a.scenes[len(a.scenes)-1]
		a.currentScene.Entered()
		log.Printf("App Pop New Scene Entered.")
	}
	if len(a.scenes) == 0 {
		log.Printf("App Quit.")
		a.db.Close()
		os.Exit(0)
	}
}

func (a *App) getTick() int {
	tm := time.Now()
	dt := tm.Nanosecond() / 1e6
	if a.lastDt == -1 {
		a.lastDt = dt
	}
	ticks := dt - a.lastDt
	if dt < a.lastDt {
		ticks = 999 - a.lastDt + dt
	}
	a.lastDt = dt
	return ticks
}
