package ui

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ui struct {
	startDt      time.Time
	fullScreen   bool
	rect         *Rect
	scenes       []Scene
	currentScene Scene
	lastDt       int
	theme        *Theme
	pref         *Preferences
	locale       *Locale
}

func init() {
	uiInstance = GetUi()
}

var uiInstance *Ui = nil

func GetUi() (a *Ui) {
	if uiInstance == nil {
		a = &Ui{
			startDt: time.Now(),
			lastDt:  -1,
			scenes:  []Scene{},
		}
		log.Printf("App init done")
	} else {
		a = uiInstance
	}
	return a
}

func (a *Ui) SetupSettings(p *Preferences) {
	a.pref = p
	log.Printf("App init preferences: %v", a.pref)
}

func (a *Ui) SetupLocale(l *Locale) {
	a.locale = l
	log.Printf("App init Locale: %v", a.locale)
}

func (a *Ui) SetupTheme(theme *Theme) {
	a.theme = theme
	log.Printf("App init theme: %v", a.theme)
}

func (a *Ui) SetupScreen(title string) {
	var w, h int
	if a.fullScreen {
		w, h = ebiten.ScreenSizeInFullscreen()
	} else {
		w, h = fitWindowSize()
	}
	ebiten.SetWindowTitle(title)
	ebiten.SetFullscreen(a.fullScreen)
	ebiten.SetWindowSize(w, h)
	a.rect = NewRect([]int{0, 0, w, h})
}

func GetLocale() *Locale {
	return GetUi().locale
}

func GetTheme() *Theme {
	return GetUi().theme
}

func GetPreferences() *Preferences {
	return GetUi().pref
}

func (a *Ui) SetFullscreen(value bool) {
	a.fullScreen = value
}

func fitWindowSize() (w int, h int) {
	ww, hh := ebiten.ScreenSizeInFullscreen()
	k := 10
	w, h = 320*k, 200*k
	for ww <= w || hh <= h {
		k -= 1
		w, h = 320*k, 200*k
	}
	return w, h
}

func (a *Ui) GetScreenSize() (w, h int) {
	return a.rect.Right(), a.rect.Bottom()
}

func (a *Ui) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		a.Pop()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF11) {
		a.ToggleFullscreen()
	}
	a.currentScene.Update(a.getTick())
	return nil
}

func (a *Ui) ToggleFullscreen() {
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
	a.rect = NewRect([]int{0, 0, w, h})
	for _, scene := range a.scenes {
		scene.Resize()
	}
	log.Println("Toggle FullScreen:", a.rect)
}

func (a *Ui) Draw(screen *ebiten.Image) {
	screen.Fill((*a.theme)["bg"])
	a.currentScene.Draw(screen)
}

func (a *Ui) Layout(oW, oH int) (int, int) {
	return oW, oH
}

func (a *Ui) Push(sc Scene) {
	a.scenes = append(a.scenes, sc)
	a.currentScene = sc
	a.currentScene.Entered()
	log.Println("Scene push")
}

func (a *Ui) Pop() {
	if len(a.scenes) > 0 {
		a.currentScene.Quit()
		idx := len(a.scenes) - 1
		a.scenes = a.scenes[:idx]
		log.Printf("App Pop Scene Quit done.")
		if len(a.scenes) == 0 {
			log.Printf("App Quit.")
			os.Exit(0)
		}
		a.currentScene = a.scenes[len(a.scenes)-1]
		a.currentScene.Entered()
		log.Printf("App Pop New Scene Entered.")
	}
}

func (a *Ui) getTick() int {
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

func (s *Ui) UpdateUpTime() string {
	durration := time.Since(s.startDt)
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
	return fmt.Sprintf("%v: %v", GetLocale().Get("lblUpTm"), result)
}
