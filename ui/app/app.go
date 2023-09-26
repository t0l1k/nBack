package app

import (
	"errors"
	"log"
	"os"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

func NewGame() *eui.Ui {
	ui := eui.GetUi()
	ui.SetTitle("Single NBack")
	ui.ApplyPreferences(LoadPreferences())
	ui.ApplyTheme(NewTheme())
	ui.ApplyLocale(NewLocale())
	return ui
}

func LoadPreferences() *eui.Preferences {
	if eui.GetPreferences() == nil {
		if _, err := os.Stat("games.db"); errors.Is(err, os.ErrNotExist) {
			log.Println("Load default settings")
			sets := NewPref()
			eui.GetUi().SetFullscreen(sets.Get("fullscreen").(bool))
			return sets
		} else {
			if sets := data.GetDb().ReadSettings(); sets == nil {
				log.Println("Load default settings")
				sets := NewPref()
				eui.GetUi().SetFullscreen(sets.Get("fullscreen").(bool))
				return sets
			} else {
				log.Println("Load saved settings", sets)
				eui.GetUi().SetFullscreen(sets.Get("fullscreen").(bool))
				return sets
			}
		}
	}
	return eui.GetPreferences()
}

func NewPref() *eui.Preferences {
	p := eui.NewPreferences()
	p.Set("fullscreen", false)
	p.Set("game type", game.Ari)
	p.Set("time to next cell", 3.0)
	p.Set("time to show cell", 0.75)
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
	p.Set("symbols count", 100)
	p.Set("ariphmetic max", 20)
	return &p
}

func NewTheme() *eui.Theme {
	theme := eui.NewTheme()
	theme.Set("bg", eui.Gray)
	theme.Set("fg", eui.White)
	theme.Set("game bg", eui.Black)
	theme.Set("game fg", eui.Gray)
	theme.Set("game active color", eui.Yellow)
	theme.Set("regular color", eui.Blue)
	theme.Set("correct color", eui.Green)
	theme.Set("warning color", eui.Orange)
	theme.Set("error color", eui.Red)
	return &theme
}

func NewLocale() *eui.Locale {
	lang := eui.GetPreferences().Get("lang")
	log.Printf("Setup Locale for %v", lang)
	switch lang {
	case "ru":
		return NewLocaleRu()
	default:
		return NewLocaleEn()
	}
}

func NewLocaleEn() *eui.Locale {
	loc := eui.NewLocale()
	loc.Set("AppName", "N-Back")
	loc.Set("btnStart", "Play")
	loc.Set("btnScore", "Score")
	loc.Set("btnOpt", "Settings")
	loc.Set("btnHelper", "Press <SPACE> start playing, <P> graph, <S> score, <F11> toggle full screen, <O> settings, <Esc> exit")
	loc.Set("btnHelperInGame", "Press <SPACE> to start playing,<F5/F6> Increase/decrease time by 0.5 seconds, <F11> toggle full screen, <Esc> exit")
	loc.Set("lblUpTm", "Up")
	loc.Set("btnPlot", "{P}")
	loc.Set("wordMax", "max")
	loc.Set("wordAvg", "average")
	loc.Set("wordGames", "Games")
	loc.Set("lblGmNr", "Game Number")
	loc.Set("lblLevel", "Level")
	loc.Set("lblDTl", "Summary")
	loc.Set("scrName", "Total for the period")
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
	loc.Set("optdeflevadv", "How many games in a row in the percent up threshold")
	loc.Set("optadv", "Percent threshold advance")
	loc.Set("optfall", "Percent threshold fallback")
	loc.Set("optgmadv", "Threshold fallback sessions")
	loc.Set("optmv", "Trials")
	loc.Set("optfc", "Factor")
	loc.Set("optexp", "Exponent")
	loc.Set("opttmnc", "Time to next cell seconds")
	loc.Set("opttmsc", "Time to show cell %")
	loc.Set("optgmtp", "Modality")
	loc.Set("optpos", "Positions")
	loc.Set("optcol", "Colors")
	loc.Set("optsym", "Symbols")
	loc.Set("optari", "Arithmetic")
	loc.Set("optmaxsym", "Max number in number game")
	loc.Set("optmaxari", "Max number in arithmetic game")
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

func NewLocaleRu() *eui.Locale {
	loc := eui.NewLocale()
	loc.Set("AppName", "Н-Назад")
	loc.Set("btnStart", "Играть")
	loc.Set("btnScore", "Итог")
	loc.Set("btnOpt", "Настройки")
	loc.Set("btnHelper", "Нажать <SPACE> начать играть,<P> график, <S> итог,<F11> на весь экран, <O> настройки, <Esc> выход")
	loc.Set("btnHelperInGame", "Нажать <SPACE> начать играть,<F5/F6> Время хода увеличить/уменьшить на 0.5 секунды, <F11> на весь экран, <Esc> выход")
	loc.Set("lblUpTm", "Прошло")
	loc.Set("btnPlot", "{Г}")
	loc.Set("wordMax", "максимально")
	loc.Set("wordAvg", "среднее")
	loc.Set("wordGames", "Игр")
	loc.Set("lblGmNr", "Номер игры")
	loc.Set("lblLevel", "Уровень")
	loc.Set("lblDTl", "Итог за день")
	loc.Set("scrName", "Итог за период")
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
	loc.Set("optdeflevadv", "Сколько игр подряд в пороге процента вверх")
	loc.Set("optadv", "Порог процента перехода вверх")
	loc.Set("optfall", "Порог процента перехода вниз")
	loc.Set("optgmadv", "Дополнительно попыток")
	loc.Set("optmv", "Ходов")
	loc.Set("optfc", "Столбик")
	loc.Set("optexp", "Степень")
	loc.Set("opttmnc", "Время до следующей ячейки в секундах")
	loc.Set("opttmsc", "Время показа ячейки %")
	loc.Set("optgmtp", "Модальность")
	loc.Set("optpos", "Позиции")
	loc.Set("optcol", "Цвета")
	loc.Set("optsym", "Символы")
	loc.Set("optari", "Арифметика")
	loc.Set("optmaxsym", "Макс число в игре с цифрами")
	loc.Set("optmaxari", "Макс число в игре арифметика")
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
