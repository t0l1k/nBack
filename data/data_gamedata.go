package data

import (
	"container/list"
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/game"
)

type GameData struct {
	GameType, DtBeg, DtEnd                                                                  string
	Id, Level, Lives, Percent, Correct, Wrong, Missed, Moves, Totalmoves, Advance, Fallback int
	Manual, Resetonerror                                                                    bool
	MovesStatus                                                                             map[int]game.Status
}

func (d *GameData) NextLevel() (int, int, string) {
	motiv := ""
	manual := eui.GetPreferences().Get("manual mode").(bool)
	adv := eui.GetPreferences().Get("threshold advance").(int)
	fall := eui.GetPreferences().Get("threshold fallback").(int)
	level := d.Level
	lives := d.Lives
	win, ok, count := GetDb().TodayData.getWinCountInManual()
	if manual {
		if !win && !ok {
			motiv = eui.GetLocale().Get("strgamemanual") + " " + eui.GetLocale().Get("strmotivdef") + "(" + strconv.Itoa(level) + ")"
			level = eui.GetPreferences().Get("default level").(int)
			lives = count
		} else if !win && ok {
			motiv = eui.GetLocale().Get("strgamemanual") + " " + eui.GetLocale().Get("strmotivmed") + "(" + strconv.Itoa(level) + ")"
			lives = count
		} else if win && ok {
			level += 1
			motiv = eui.GetLocale().Get("strgamemanual") + " " + eui.GetLocale().Get("strmotivup") + "(" + strconv.Itoa(level) + ")"
			lives = 0
		}
	} else {
		if lives == 0 || count > 0 {
			level = eui.GetPreferences().Get("default level").(int)
			lives = eui.GetPreferences().Get("threshold fallback sessions").(int)
		} else if d.Percent >= adv {
			level += 1
			lives = eui.GetPreferences().Get("threshold fallback sessions").(int)
			motiv = eui.GetLocale().Get("strgameclassic") + " " + eui.GetLocale().Get("strmotivup") + "(" + strconv.Itoa(level) + ")"
		} else if d.Percent >= fall && d.Percent < adv {
			motiv = eui.GetLocale().Get("strgameclassic") + " " + eui.GetLocale().Get("strmotivmed") + "(" + strconv.Itoa(level) + ")"
		} else if d.Percent < fall {
			if lives == 1 {
				if level > 1 {
					level -= 1
					lives = eui.GetPreferences().Get("threshold fallback sessions").(int)
				}
				motiv = eui.GetLocale().Get("strgameclassic") + " " + eui.GetLocale().Get("strmotivdwn") + "(" + strconv.Itoa(level) + ")"
			} else if lives > 1 {
				lives -= 1
				motiv = eui.GetLocale().Get("strgameclassic") + " " + eui.GetLocale().Get("strmotivadv") + "(" + strconv.Itoa(level) + ")"
			}
		}
	}
	return level, lives, motiv
}
func (d GameData) MovesColor() (moves, colors list.List) {
	theme := eui.GetTheme()
	colorNil := theme.Get("game fg")
	colorRegular := theme.Get("regular color")
	colorCorrect := theme.Get("correct color")
	colorError := theme.Get("error color")
	colorWarning := theme.Get("warning color")
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
	theme := eui.GetTheme()
	colorRegular := theme.Get("regular color")
	colorCorrect := theme.Get("correct color")
	colorError := theme.Get("error color")
	colorWarning := theme.Get("warning color")
	adv := eui.GetPreferences().Get("threshold advance").(int)
	fall := eui.GetPreferences().Get("threshold fallback").(int)
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
	ss := fmt.Sprintf(" %vB%v %v%% %v:%v (%v:%v %v:%v)%v %v:%v [%v]",
		q.GameType,
		q.Level,
		q.Percent,
		eui.GetLocale().Get("wordrgt"),
		q.Correct,
		eui.GetLocale().Get("worderr"),
		q.Wrong,
		eui.GetLocale().Get("wordmissed"),
		q.Missed,
		q.Wrong+q.Missed,
		eui.GetLocale().Get("wordmove"),
		q.Moves,
		dStr)
	if eui.GetPreferences().Get("reset on first wrong").(bool) {
		ss = fmt.Sprintf(" %vB%v %v%% %v:%v (%v:%v %v:%v)%v %v:(%v/%v) [%v]",
			q.GameType,
			q.Level,
			q.Percent,
			eui.GetLocale().Get("wordrgt"),
			q.Correct,
			eui.GetLocale().Get("worderr"),
			q.Wrong,
			eui.GetLocale().Get("wordmissed"),
			q.Missed,
			q.Wrong+q.Missed,
			eui.GetLocale().Get("wordmove"),
			q.Moves,
			q.Totalmoves,
			dStr)
	}
	return ss
}
