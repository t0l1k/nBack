package ui

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var uiInstance *Ui = nil

func init() {
	uiInstance = GetUi()
}

func GetUi() (a *Ui) {
	if uiInstance == nil {
		a = &Ui{
			startDt:      time.Now(),
			lastDt:       -1,
			scenes:       []Scene{},
			notification: nil,
		}
		log.Printf("App init done")
	} else {
		a = uiInstance
	}
	return a
}

func Init(f *Ui) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ebiten.SetWindowTitle(f.title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	var w, h int
	if GetUi().fullScreen {
		w, h = ebiten.ScreenSizeInFullscreen()
	} else {
		w, h = fitWindowSize()
	}
	ebiten.SetWindowSize(w, h)
	ebiten.SetFullscreen(GetUi().fullScreen)
	GetUi().rect = NewRect([]int{0, 0, w, h})
}

func Run(sc Scene) {
	Push(sc)
	if err := ebiten.RunGame(GetUi()); err != nil {
		log.Fatal(err)
	}
}

func Quit() {}

func GetLocale() *Locale {
	return GetUi().locale
}

func GetTheme() *Theme {
	return GetUi().theme
}

func GetPreferences() *Preferences {
	return GetUi().pref
}

func Push(sc Scene) {
	GetUi().scenes = append(GetUi().scenes, sc)
	GetUi().currentScene = sc
	GetUi().currentScene.Entered()
	log.Println("Scene push")
}

func Pop() {
	if len(GetUi().scenes) > 0 {
		GetUi().currentScene.Quit()
		idx := len(GetUi().scenes) - 1
		GetUi().scenes = GetUi().scenes[:idx]
		log.Printf("App Pop Scene Quit done.")
		if len(GetUi().scenes) == 0 {
			log.Printf("App Quit.")
			os.Exit(0)
		}
		GetUi().currentScene = GetUi().scenes[len(GetUi().scenes)-1]
		GetUi().currentScene.Entered()
		log.Printf("App Pop New Scene Entered.")
	}
}

func fitWindowSize() (w int, h int) {
	ww, hh := ebiten.ScreenSizeInFullscreen()
	k := 10
	w, h = 320*k, 180*k
	for ww <= w*2 || hh <= h {
		k -= 1
		w, h = 320*k, 200*k
	}
	return w, h
}
