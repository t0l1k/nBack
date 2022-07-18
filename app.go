package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type App struct {
	startDt      time.Time
	fullScreen   bool
	rect         *ui.Rect
	scenes       []ui.Scene
	currentScene ui.Scene
	lastDt       int
	db           *Db
	theme        *Theme
	preferences  *Setting
}

var app *App

func NewGame() *App {
	app = getApp()
	app.Push(NewSceneToday())
	return app
}

func getApp() (a *App) {
	if app == nil {
		db := &Db{}
		db.Setup()
		pref := NewSettings()
		db.ReadSettings(pref)
		var w, h int
		if pref.FullScreen {
			w, h = ebiten.ScreenSizeInFullscreen()
		} else {
			w, h = fitWindowSize()
		}
		ebiten.SetWindowTitle("nBack")
		ebiten.SetFullscreen(pref.FullScreen)
		ebiten.SetWindowSize(w, h)
		rect := ui.NewRect([]int{0, 0, w, h})
		scns := []ui.Scene{}
		a = &App{
			startDt:     time.Now(),
			fullScreen:  pref.FullScreen,
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

func getDb() *Db {
	return getApp().db
}

func getTheme() *Theme {
	return getApp().theme
}

func getPreferences() *Setting {
	return getApp().preferences
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

func (a *App) GetScreenSize() (w, h int) {
	return a.rect.Right(), a.rect.Bottom()
}

func (a *App) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		a.Pop()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF11) {
		a.toggleFullscreen()
	}
	a.currentScene.Update(a.getTick())
	return nil
}

func (a *App) toggleFullscreen() {
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
	log.Println("Toggle FullScreen:", a.rect)
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(a.theme.Bg)
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

func (s *App) updateUpTime() string {
	durration := time.Since(getApp().startDt)
	d := durration.Round(time.Second)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	sec := d / time.Second
	result := ""
	if hours > 0 {
		result = fmt.Sprintf("%02v:%02v:%02v", int(hours), int(minutes), int(sec))
	} else {
		result = fmt.Sprintf("%02v:%02v", int(minutes), int(sec))
	}
	return fmt.Sprintf("up: %v", result)
}
