package data

import (
	"container/list"
	"fmt"
	"image/color"
	"time"

	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui"
)

type GameData struct {
	GameType, DtBeg, DtEnd                                                                  string
	Id, Level, Lives, Percent, Correct, Wrong, Missed, Moves, Totalmoves, Advance, Fallback int
	Manual, Resetonerror                                                                    bool
	MovesStatus                                                                             map[int]game.Status
}

func (d *GameData) NextLevel() (int, int, string) {
	motiv := ""
	manual := (*ui.GetPreferences())["manual mode"].(bool)
	adv := (*ui.GetPreferences())["threshold advance"].(int)
	fall := (*ui.GetPreferences())["threshold fallback"].(int)
	level := d.Level
	lives := d.Lives
	if manual {
		win, ok, count := GetDb().TodayData.getWinCountInManual()
		if !win && !ok {
			motiv = ui.GetLocale().Get("strgamemanual") + " " + ui.GetLocale().Get("strmotivdef")
			level = (*ui.GetPreferences())["default level"].(int)
			lives = count
		} else if !win && ok {
			motiv = ui.GetLocale().Get("strgamemanual") + " " + ui.GetLocale().Get("strmotivmed")
			lives = count
		} else if win && ok {
			motiv = ui.GetLocale().Get("strgamemanual") + " " + ui.GetLocale().Get("strmotivup")
			level += 1
			lives = 0
		}
	} else if d.Percent >= adv {
		level += 1
		lives = (*ui.GetPreferences())["threshold fallback sessions"].(int)
		motiv = ui.GetLocale().Get("strgameclassic") + " " + ui.GetLocale().Get("strmotivup")
	} else if d.Percent >= fall && d.Percent < adv {
		motiv = ui.GetLocale().Get("strgameclassic") + " " + ui.GetLocale().Get("strmotivmed")
	} else if d.Percent < fall {
		if lives == 1 {
			motiv = ui.GetLocale().Get("strgameclassic") + " " + ui.GetLocale().Get("strmotivdwn")
			if level > 1 {
				level -= 1
				lives = (*ui.GetPreferences())["threshold fallback sessions"].(int)
			}
		} else if lives > 1 {
			motiv = ui.GetLocale().Get("strgameclassic") + " " + ui.GetLocale().Get("strmotivadv")
			lives -= 1
		}
	}
	return level, lives, motiv
}
func (d GameData) MovesColor() (moves, colors list.List) {
	theme := ui.GetTheme()
	colorNil := (*theme)["game fg"]
	colorRegular := (*theme)["regular color"]
	colorCorrect := (*theme)["correct color"]
	colorError := (*theme)["error color"]
	colorWarning := (*theme)["warning color"]
	for k, v := range d.MovesStatus {
		moves.PushBack(k)
		clr := colorNil
		switch v {
		case game.Neutral, game.Regular:
			clr = colorRegular
		case game.Correct:
			clr = colorCorrect
		case game.Error:
			clr = colorError
		case game.Warning:
			clr = colorWarning
		}
		colors.PushBack(clr)
	}
	if d.Moves < d.Totalmoves {
		for i := d.Moves + 1; i < d.Totalmoves+1; i++ {
			moves.PushBack(i)
			colors.PushBack(colorNil)
		}
	}
	return moves, colors
}

func (d GameData) BgColor() (result color.Color) {
	theme := ui.GetTheme()
	colorRegular := (*theme)["regular color"]
	colorCorrect := (*theme)["correct color"]
	colorError := (*theme)["error color"]
	colorWarning := (*theme)["warning color"]
	adv := (*ui.GetPreferences())["threshold advance"].(int)
	fall := (*ui.GetPreferences())["threshold fallback"].(int)
	if d.Percent >= adv {
		result = colorRegular
	} else if d.Percent >= fall && d.Percent < adv {
		result = colorCorrect
	} else if d.Percent < fall {
		if d.Lives <= 1 {
			result = colorError
		} else if d.Lives > 1 {
			result = colorWarning
		}
	}
	return
}

func (q GameData) ShortStr() string {
	return fmt.Sprintf("%vB%v %v%% ", q.GameType, q.Level, q.Percent)

}
func (q GameData) String() string {
	var durration time.Duration
	dtFormat := "2006-01-02 15:04:05.000"
	dtBeg, err := time.Parse(dtFormat, q.DtBeg)
	if err != nil {
		panic(err)
	}
	dtEnd, err := time.Parse(dtFormat, q.DtEnd)
	if err != nil {
		panic(err)
	}
	durration = dtEnd.Sub(dtBeg)
	mSec := durration.Milliseconds() / 1e3
	sec := durration.Seconds()
	m := int(sec / 60)
	seconds := int(sec) % 60
	dStr := fmt.Sprintf("%02v:%02v.%03v", m, seconds, int(mSec))
	ss := fmt.Sprintf("#%v %vB%v %v%% %v:%v %v:%v %v:%v %v:%v [%v]",
		GetDb().TodayGamesCount,
		q.GameType,
		q.Level,
		q.Percent,
		ui.GetLocale().Get("wordrgt"),
		q.Correct,
		ui.GetLocale().Get("worderr"),
		q.Wrong,
		ui.GetLocale().Get("wordmissed"),
		q.Missed,
		ui.GetLocale().Get("wordmove"),
		q.Moves,
		dStr)
	if (*ui.GetPreferences())["reset on first wrong"].(bool) {
		ss = fmt.Sprintf("#%v %vB%v %v%% %v:%v %v:%v %v:%v %v:(%v/%v) [%v]",
			GetDb().TodayGamesCount,
			q.GameType,
			q.Level,
			q.Percent,
			ui.GetLocale().Get("wordrgt"),
			q.Correct,
			ui.GetLocale().Get("worderr"),
			q.Wrong,
			ui.GetLocale().Get("wordmissed"),
			q.Missed,
			ui.GetLocale().Get("wordmove"),
			q.Moves,
			q.Totalmoves,
			dStr)
	}
	return ss
}
