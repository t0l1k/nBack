package game

import (
	"errors"
	"image/color"
	"log"
	"os"

	"github.com/t0l1k/nBack/ui"
)

const (
	pos string = "p"
	col string = "c"
	sym string = "s"
)

func NewGame() *ui.App {
	app := ui.GetApp()
	app.SetupSettings(LoadPreferences())
	app.SetupTheme(NewTheme())
	app.SetupScreen("Single nBack")
	app.Push(NewSceneToday())
	return app
}

func ApplyPreferences(value *ui.Preferences) *ui.Preferences {
	ui.GetApp().SetupSettings(value)
	return ui.GetPreferences()
}

func LoadPreferences() *ui.Preferences {
	if ui.GetPreferences() == nil {
		if _, err := os.Stat("games.db"); errors.Is(err, os.ErrNotExist) {
			log.Println("Load default settings")
			sets := NewPref()
			ui.GetApp().SetFullscreen(sets.Get("fullscreen").(bool))
			return sets
		} else {
			if sets := getDb().ReadSettings(); sets == nil {
				log.Println("Load default settings")
				sets := NewPref()
				ui.GetApp().SetFullscreen(sets.Get("fullscreen").(bool))
				return sets
			} else {
				log.Println("Load saved settings", sets)
				ui.GetApp().SetFullscreen(sets.Get("fullscreen").(bool))
				return sets
			}
		}
	}
	return ui.GetPreferences()
}

func NewPref() *ui.Preferences {
	p := ui.NewPreferences()
	p["game type"] = sym
	p["symbols count"] = 99
	p["time to next cell"] = 2.0
	p["time to show cell"] = 0.5
	p["trials"] = 5 //20 classic = trials+factor*level**exponent
	p["trials factor"] = 1
	p["trials exponent"] = 2
	p["threshold advance"] = 80
	p["threshold fallback"] = 50
	p["threshold fallback sessions"] = 3
	p["default level"] = 1 // Level in manul mode and first game level today
	p["manual mode"] = false
	p["manual advance"] = 3 // games with 100% to next level in manual mode, 0 same level
	p["reset on first wrong"] = true
	p["random repition"] = 12.5 // Random Repition
	p["use center cell"] = false
	p["show grid"] = true
	p["show crosshair"] = true
	p["feedback on user move"] = true
	p["grid size"] = 3
	p["pause to rest"] = 5
	p["fullscreen"] = false
	return &p
}

var (
	black  = color.RGBA{0, 0, 0, 255}
	gray   = color.RGBA{128, 128, 128, 255}
	silver = color.RGBA{192, 192, 192, 255}
	white  = color.RGBA{255, 255, 255, 255}

	orange  = color.RGBA{255, 165, 0, 255}
	fuchsia = color.RGBA{255, 0, 255, 255}
	purple  = color.RGBA{128, 0, 128, 255}
	red     = color.RGBA{255, 0, 0, 255}
	maroon  = color.RGBA{128, 0, 0, 255}

	yellow      = color.RGBA{255, 255, 0, 255}
	greenYellow = color.RGBA{173, 255, 47, 255}
	yellowGreen = color.RGBA{154, 205, 50, 255}
	olive       = color.RGBA{128, 128, 0, 255}
	lime        = color.RGBA{0, 255, 0, 255}
	green       = color.RGBA{0, 128, 0, 255}

	aqua = color.RGBA{0, 255, 255, 255}
	teal = color.RGBA{0, 128, 128, 255}
	blue = color.RGBA{0, 0, 255, 255}
	navy = color.RGBA{0, 0, 128, 255}
)
var colors = []color.Color{blue, aqua, green, olive, yellow, red, purple, orange, white, gray}

// var colors = []color.Color{navy, blue, teal, aqua, green, lime, olive, yellowGreen, greenYellow, yellow, maroon, red, purple, fuchsia, orange, white, silver, gray}

func NewTheme() *ui.Theme {
	theme := ui.NewTheme()
	theme.Set("bg", gray)
	theme.Set("fg", white)
	theme.Set("game bg", black)
	theme.Set("game fg", gray)
	theme.Set("game active color", yellow)
	theme.Set("regular color", blue)
	theme.Set("correct color", green)
	theme.Set("warning color", orange)
	theme.Set("error color", red)
	return &theme
}
