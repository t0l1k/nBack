package game

import (
	"fmt"
	"image/color"
	"time"
)

type GameData struct {
	dtBeg, dtEnd                                                                    string
	id, level, lives, percent, correct, wrong, moves, totalmoves, advance, fallback int
	manual, resetonerror                                                            bool
}

func (d *GameData) NextLevel() (int, int, string) {
	motiv := ""
	manual := getPreferences().Manual
	adv := getPreferences().ThresholdAdvance
	fall := getPreferences().ThresholdFallback
	level := d.level
	lives := d.lives
	if manual {
		win, ok, count := getDb().todayData.getWinCountInManual()
		if !win && !ok {
			motiv = "Manual game mode. Level default."
			level = getPreferences().DefaultLevel
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
		lives = getPreferences().ThresholdFallbackSessions
		motiv = "Classic game mode. Excellent result! Level up!"
	} else if d.percent >= fall && d.percent < adv {
		motiv = "Classic game mode. Good result! One more time this level!"
	} else if d.percent < fall {
		if lives == 1 {
			motiv = "Classic game mode. Let's improve the results! Level down!"
			if level > 1 {
				level -= 1
				lives = getPreferences().ThresholdFallbackSessions
			}
		} else if lives > 1 {
			motiv = "Classic game mode. Let's improve the results! Let's have an extra try!"
			lives -= 1
		}
	}
	return level, lives, motiv
}

func (d GameData) BgColor() (result color.Color) {
	theme := getTheme()
	colorRegular := theme.RegularColor
	colorCorrect := theme.CorrectColor
	colorError := theme.ErrorColor
	colorWarning := theme.WarningColor
	adv := getPreferences().ThresholdAdvance
	fall := getPreferences().ThresholdFallback
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
	return fmt.Sprintf("nB%v %v%% ", q.level, q.percent)

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
	ss := fmt.Sprintf("#%v nB%v %v%% correct:%v wrong:%v moves:%v [%v]",
		getDb().todayGamesCount,
		q.level,
		q.percent,
		q.correct,
		q.wrong,
		q.moves,
		dStr)
	if getPreferences().ResetOnFirstWrong {
		ss = fmt.Sprintf("#%v nB%v %v%% correct:%v wrong:%v moves:(%v/%v) [%v]",
			getDb().todayGamesCount,
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
