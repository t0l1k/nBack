package app

import (
	"errors"
	"log"
	"os"

	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui"
)

func NewGame() *ui.Ui {
	ui := ui.GetUi()
	ui.SetTitle("Single NBack")
	ui.ApplyPreferences(LoadPreferences())
	ui.ApplyTheme(NewTheme())
	ui.ApplyLocale(NewLocale())
	return ui
}

func LoadPreferences() *ui.Preferences {
	if ui.GetPreferences() == nil {
		if _, err := os.Stat("games.db"); errors.Is(err, os.ErrNotExist) {
			log.Println("Load default settings")
			sets := NewPref()
			ui.GetUi().SetFullscreen(sets.Get("fullscreen").(bool))
			return sets
		} else {
			if sets := data.GetDb().ReadSettings(); sets == nil {
				log.Println("Load default settings")
				sets := NewPref()
				ui.GetUi().SetFullscreen(sets.Get("fullscreen").(bool))
				return sets
			} else {
				log.Println("Load saved settings", sets)
				ui.GetUi().SetFullscreen(sets.Get("fullscreen").(bool))
				return sets
			}
		}
	}
	return ui.GetPreferences()
}

func NewPref() *ui.Preferences {
	p := ui.NewPreferences()
	p.Set("fullscreen", false)
	p.Set("game type", game.Sym)
	p.Set("symbols count", 99)
	p.Set("time to next cell", 2.5)
	p.Set("time to show cell", 0.5)
	p.Set("trials", 20) //20 classic = trials+factor*level**exponent
	p.Set("trials factor", 1)
	p.Set("trials exponent", 2)
	p.Set("threshold advance", 80)
	p.Set("threshold fallback", 50)
	p.Set("threshold fallback sessions", 3)
	p.Set("default level", 1) // Level in manul mode and first game level today
	p.Set("manual mode", false)
	p.Set("manual advance", 3) // games with 100% to next level in manual mode, 0 same level
	p.Set("reset on first wrong", false)
	p.Set("random repition", 12.5) // Random Repition
	p.Set("use center cell", false)
	p.Set("show grid", true)
	p.Set("show crosshair", true)
	p.Set("feedback on user move", true)
	p.Set("grid size", 3)
	p.Set("pause to rest", 5)
	p.Set("lang", "ru")
	return &p
}

func NewTheme() *ui.Theme {
	theme := ui.NewTheme()
	theme.Set("bg", ui.Gray)
	theme.Set("fg", ui.White)
	theme.Set("game bg", ui.Black)
	theme.Set("game fg", ui.Gray)
	theme.Set("game active color", ui.Yellow)
	theme.Set("regular color", ui.Blue)
	theme.Set("correct color", ui.Green)
	theme.Set("warning color", ui.Orange)
	theme.Set("error color", ui.Red)
	return &theme
}

func NewLocale() *ui.Locale {
	lang := ui.GetPreferences().Get("lang")
	log.Printf("Setup Locale for %v", lang)
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
	loc.Set("strmotivadv", "Let's get better results! Level repeat!")
	loc.Set("wordrgt", "correct")
	loc.Set("worderr", "mistakes")
	loc.Set("wordmissed", "omitted")
	loc.Set("wordmove", "moves")
	loc.Set("wordGame", "Game")
	loc.Set("wordstepback", "steps back")
	loc.Set("wordhand", "manual mode")
	loc.Set("wordcclassic", "classic mode")
	loc.Set("wordnewsess", "New Session")
	loc.Set("notifhere", "Notifications here")
	loc.Set("inc", "Increase")
	loc.Set("dec", "Decrease")
	loc.Set("by", "by")
	loc.Set("sec", "seconds")

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
	loc.Set("strmotivadv", "Дополнительная попытка! Уровень повторим!")
	loc.Set("wordrgt", "правильных")
	loc.Set("worderr", "ошибок")
	loc.Set("wordmissed", "пропущеных")
	loc.Set("wordmove", "ходов")
	loc.Set("wordGame", "Игра")
	loc.Set("wordstepback", "шага назад")
	loc.Set("wordhand", "режим на ручнике")
	loc.Set("wordcclassic", "режим классика")
	loc.Set("wordnewsess", "Новая сессия")
	loc.Set("notifhere", "Тут уведомления")
	loc.Set("inc", "Увеличили")
	loc.Set("dec", "Уменьшили")
	loc.Set("by", "на")
	loc.Set("sec", "секунд")

	return &loc
}
