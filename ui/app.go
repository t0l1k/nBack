package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ui struct {
	title        string
	startDt      time.Time
	fullScreen   bool
	rect, last   *Rect
	scenes       []Scene
	currentScene Scene
	lastDt       int
	theme        *Theme
	pref         *Preferences
	locale       *Locale
	notification *Notification
}

func (a *Ui) SetTitle(title string) {
	a.title = title
}

func (a *Ui) ApplyPreferences(value *Preferences) *Preferences {
	a.pref = value
	log.Printf("App init preferences.")
	return a.pref
}

func (a *Ui) ApplyLocale(l *Locale) {
	a.locale = l
	log.Printf("App init Locale.")
}

func (a *Ui) ApplyTheme(theme *Theme) {
	a.theme = theme
	log.Printf("App init theme")
}

func (a *Ui) ShowNotification(text string) {
	w, h := int(float64(a.rect.W)*0.50), int(float64(a.rect.H)*0.05)
	x, y := a.rect.CenterX()-w/2, 0
	rect := []int{x, y, w, h}
	bg := GetTheme().Get("bg")
	fg := GetTheme().Get("fg")
	a.notification = NewNotification(text, 2, rect, bg, fg)
	log.Printf("Show Notification %v", text)
}

func (a *Ui) setRect(w int, h int) {
	a.last = a.rect
	a.rect = NewRect([]int{0, 0, w, h})
}

func (a *Ui) SetFullscreen(value bool) {
	a.fullScreen = value
}

func (a *Ui) GetScreenSize() (w, h int) {
	return a.rect.Right(), a.rect.Bottom()
}

func (a *Ui) Layout(w, h int) (int, int) {
	return w, h
}

func (a *Ui) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		Pop()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF11) {
		a.ToggleFullscreen()
	}
	var w, h int
	if a.fullScreen {
		w, h = ebiten.ScreenSizeInFullscreen()
	} else {
		w, h = ebiten.WindowSize()
	}
	w1, h1 := a.GetScreenSize()
	if w != w1 || h != h1 {
		a.setRect(w, h)
		for _, scene := range a.scenes {
			scene.Resize()
		}
		log.Printf("Resized: %v %v %v %v", w, h, w1, h1)
	}
	tick := a.getTick()
	a.currentScene.Update(tick)
	if a.notification != nil {
		a.notification.Update(tick)
		if !a.notification.Show {
			a.notification = nil
			log.Printf("Notification off")
		}
	}
	return nil
}

func (a *Ui) Draw(screen *ebiten.Image) {
	screen.Fill(a.theme.Get("bg"))
	a.currentScene.Draw(screen)
	if a.notification != nil {
		a.notification.Draw(screen)
	}
}

func (a *Ui) ToggleFullscreen() {
	a.fullScreen = !a.fullScreen
	var w, h int
	if a.fullScreen {
		ebiten.SetFullscreen(a.fullScreen)
		w, h = ebiten.ScreenSizeInFullscreen()
	} else {
		if a.last == nil {
			w, h := ebiten.WindowSize()
			a.last = NewRect([]int{0, 0, w, h})
		}
		w, h = a.last.W, a.last.H
	}
	ebiten.SetWindowSize(w, h)
	ebiten.SetFullscreen(a.fullScreen)
	a.setRect(w, h)
	for _, scene := range a.scenes {
		scene.Resize()
	}
	log.Println("Toggle FullScreen:", a.rect)
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
