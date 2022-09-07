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
	app.SetupLocale(NewLocale())
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
	p.Set("lang", "ru")
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

func NewLocale() *ui.Locale {
	lang := ui.GetPreferences().Get("lang")
	log.Printf("Setup Locale %v", lang)
	switch lang {
	case "ru":
		return NewLocaleRu()
	default:
		return NewLocaleEn()
	}
}

func NewLocaleEn() *ui.Locale {
	loc := ui.NewLocale()
	loc.Set("AppName", "N-Back")
	loc.Set("btnStart", "Play")
	loc.Set("btnScore", "Score")
	loc.Set("btnOpt", "Settings")
	loc.Set("btnHelper", "Press <SPACE> start playing, <P> graph, <S> score, <F11> toggle full screen, <O> settings, <Esc> exit")
	loc.Set("lblUpTm", "Up")
	loc.Set("btnPlot", "{P}")
	loc.Set("wordMax", "max")
	loc.Set("wordAvg", "average")
	loc.Set("wordGames", "Games")
	loc.Set("lblGmNr", "Game Number")
	loc.Set("lblLevel", "Level")
	loc.Set("lblDTl", "Summary")
	loc.Set("scrName", "Games for the period")
	loc.Set("scrResultNil", "There is no result yet to show.")
	loc.Set("scrResultTtl", "Total games")
	loc.Set("lblDays", "Days")
	loc.Set("btnReset", "Reset")
	loc.Set("btnSave", "Save")
	loc.Set("optfs", "Run in full screen")
	loc.Set("optcc", "Use cell in center")
	loc.Set("optfeedback", "Feedback on user move")
	loc.Set("optreset", "Reset on first error")
	loc.Set("optmanual", "Manual mode")
	loc.Set("optgrid", "Show grid")
	loc.Set("optcross", "Show crosshair")
	loc.Set("optgridsz", "Grid Size")
	loc.Set("optrr", "Random repition")
	loc.Set("optpause", "Compulsory pause after the game")
	loc.Set("optdeflev", "Default level")
	loc.Set("optdeflevadv", "Games 100% transition")
	loc.Set("optadv", "Threshold advance")
	loc.Set("optfall", "Threshold fallback")
	loc.Set("optgmadv", "Threshold fallback sessions")
	loc.Set("optmv", "Trials")
	loc.Set("optfc", "Factor")
	loc.Set("optexp", "Exponent")
	loc.Set("opttmnc", "Time to next cell")
	loc.Set("opttmsc", "Time to show cell")
	loc.Set("optgmtp", "Game Type")
	loc.Set("optpos", "Positions")
	loc.Set("optcol", "Colors")
	loc.Set("optsym", "Symbols")
	loc.Set("optlang", "Language")
	loc.Set("strgamemanual", "Manual game mode.")
	loc.Set("strgameclassic", "Classic game mode.")
	loc.Set("strmotivdef", "Default level.")
	loc.Set("strmotivmed", "Good result! Once again this level!")
	loc.Set("strmotivup", "Great result! Level up!")
	loc.Set("strmotivdwn", "Let's improve even more! Level down!")
	loc.Set("strmotivadv", "Let's get better results! Extra try!")
	loc.Set("wordrgt", "correct")
	loc.Set("worderr", "mistakes")
	loc.Set("wordmissed", "omitted")
	loc.Set("wordmove", "moves")
	loc.Set("wordGame", "Game")
	loc.Set("wordstepback", "steps back")
	loc.Set("wordhand", "manual mode")
	loc.Set("wordcclassic", "classic mode")
	loc.Set("wordnewsess", "New Session")

	return &loc
}

func NewLocaleRu() *ui.Locale {
	loc := ui.NewLocale()
	loc.Set("AppName", "Н-Назад")
	loc.Set("btnStart", "Играть")
	loc.Set("btnScore", "Итог")
	loc.Set("btnOpt", "Настройки")
	loc.Set("btnHelper", "Нажать <SPACE> начать играть,<P> график, <S> итог,<F11> на весь экран, <O> настройки, <Esc> выход")
	loc.Set("lblUpTm", "Прошло")
	loc.Set("btnPlot", "{Г}")
	loc.Set("wordMax", "максимально")
	loc.Set("wordAvg", "среднее")
	loc.Set("wordGames", "Игр")
	loc.Set("lblGmNr", "Номер игры")
	loc.Set("lblLevel", "Уровень")
	loc.Set("lblDTl", "Итог за день")
	loc.Set("scrName", "Игры за период")
	loc.Set("scrResultNil", "Ещё нет результата, что показать.")
	loc.Set("scrResultTtl", "Всего игр")
	loc.Set("lblDays", "Дней")
	loc.Set("btnReset", "Обнулить")
	loc.Set("btnSave", "Сохранить")
	loc.Set("optfs", "Запуск на весь экран")
	loc.Set("optcc", "Использовать ячейку в центре")
	loc.Set("optfeedback", "Отклик на ход игры")
	loc.Set("optreset", "Сброс при первой ошибке")
	loc.Set("optmanual", "Игра на ручнике")
	loc.Set("optgrid", "Показать сетку")
	loc.Set("optcross", "Показать прицел")
	loc.Set("optgridsz", "Размер сетки")
	loc.Set("optrr", "Процент повторов")
	loc.Set("optpause", "Обязательная пауза после игры")
	loc.Set("optdeflev", "Уровень по умолчанию")
	loc.Set("optdeflevadv", "Игр на 100% переход")
	loc.Set("optadv", "Процент перехода вверх")
	loc.Set("optfall", "Процент перехода вниз")
	loc.Set("optgmadv", "Дополнительно попыток")
	loc.Set("optmv", "Ходов")
	loc.Set("optfc", "Столбик")
	loc.Set("optexp", "Степень")
	loc.Set("opttmnc", "Время до следующей ячейки")
	loc.Set("opttmsc", "Время показа ячейки")
	loc.Set("optgmtp", "Тип игры")
	loc.Set("optpos", "Позиции")
	loc.Set("optcol", "Цвета")
	loc.Set("optsym", "Символы")
	loc.Set("optlang", "Язык")
	loc.Set("strgamemanual", "Режим игры на ручнике.")
	loc.Set("strgameclassic", "Режим игры классика.")
	loc.Set("strmotivdef", "Уровень по умолчанию.")
	loc.Set("strmotivmed", "Хороший результат! Еще раз этот уровень!")
	loc.Set("strmotivup", "Отличный результат! Уровень повышен!")
	loc.Set("strmotivdwn", "Улучшим ещё результаты! Уровень вниз!")
	loc.Set("strmotivadv", "Улучшим ещё результаты! Дополнительная попытка!")
	loc.Set("wordrgt", "правильных")
	loc.Set("worderr", "ошибок")
	loc.Set("wordmissed", "пропущеных")
	loc.Set("wordmove", "ходов")
	loc.Set("wordGame", "Игра")
	loc.Set("wordstepback", "шага назад")
	loc.Set("wordhand", "режим на ручнике")
	loc.Set("wordcclassic", "режим классика")
	loc.Set("wordnewsess", "Новая сессия")

	return &loc
}
