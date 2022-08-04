package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/t0l1k/nBack/ui"
)

type GameData struct {
	gameType, dtBeg, dtEnd                                                          string
	id, level, lives, percent, correct, wrong, moves, totalmoves, advance, fallback int
	manual, resetonerror                                                            bool
}

func (d *GameData) NextLevel() (int, int, string) {
	motiv := ""
	manual := (*ui.GetPreferences())["manual mode"].(bool)
	adv := (*ui.GetPreferences())["threshold advance"].(int)
	fall := (*ui.GetPreferences())["threshold fallback"].(int)
	level := d.level
	lives := d.lives
	if manual {
		win, ok, count := getDb().todayData.getWinCountInManual()
		if !win && !ok {
			motiv = "Manual game mode. Level default."
			level = (*ui.GetPreferences())["default level"].(int)
			lives = count
		} else if !win && ok {
			motiv = "Manual game mode. Good result! One more time this level!"
			lives = count
		} else if win && ok {
			motiv = "Manual game mode. Excellent result! Level up!"
			level += 1
			lives = 0
		}
	} else if d.percent >= adv {
		level += 1
		lives = (*ui.GetPreferences())["threshold fallback sessions"].(int)
		motiv = "Classic game mode. Excellent result! Level up!"
	} else if d.percent >= fall && d.percent < adv {
		motiv = "Classic game mode. Good result! One more time this level!"
	} else if d.percent < fall {
		if lives == 1 {
			motiv = "Classic game mode. Let's improve the results! Level down!"
			if level > 1 {
				level -= 1
				lives = (*ui.GetPreferences())["threshold fallback sessions"].(int)
			}
		} else if lives > 1 {
			motiv = "Classic game mode. Let's improve the results! Let's have an extra try!"
			lives -= 1
		}
	}
	return level, lives, motiv
}

func (d GameData) BgColor() (result color.Color) {
	theme := ui.GetTheme()
	colorRegular := (*theme)["regular color"]
	colorCorrect := (*theme)["correct color"]
	colorError := (*theme)["error color"]
	colorWarning := (*theme)["warning color"]
	adv := (*ui.GetPreferences())["threshold advance"].(int)
	fall := (*ui.GetPreferences())["threshold fallback"].(int)
	if d.percent >= adv {
		result = colorRegular
	} else if d.percent >= fall && d.percent < adv {
		result = colorCorrect
	} else if d.percent < fall {
		if d.lives <= 1 {
			result = colorError
		} else if d.lives > 1 {
			result = colorWarning
		}
	}
	return
}

func (q GameData) ShortStr() string {
	return fmt.Sprintf("%vB%v %v%% ", q.gameType, q.level, q.percent)

}
func (q GameData) String() string {
	var durration time.Duration
	dtFormat := "2006-01-02 15:04:05.000"
	dtBeg, err := time.Parse(dtFormat, q.dtBeg)
	if err != nil {
		panic(err)
	}
	dtEnd, err := time.Parse(dtFormat, q.dtEnd)
	if err != nil {
		panic(err)
	}
	durration = dtEnd.Sub(dtBeg)
	mSec := durration.Milliseconds() / 1e3
	sec := durration.Seconds()
	m := int(sec / 60)
	seconds := int(sec) % 60
	dStr := fmt.Sprintf("%02v:%02v.%03v", m, seconds, int(mSec))
	ss := fmt.Sprintf("#%v %vB%v %v%% correct:%v wrong:%v moves:%v [%v]",
		getDb().todayGamesCount,
		q.gameType,
		q.level,
		q.percent,
		q.correct,
		q.wrong,
		q.moves,
		dStr)
	if (*ui.GetPreferences())["reset on first wrong"].(bool) {
		ss = fmt.Sprintf("#%v %vB%v %v%% correct:%v wrong:%v moves:(%v/%v) [%v]",
			getDb().todayGamesCount,
			q.gameType,
			q.level,
			q.percent,
			q.correct,
			q.wrong,
			q.moves,
			q.totalmoves,
			dStr)
	}
	return ss
}
