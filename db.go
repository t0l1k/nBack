package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"image/color"
	"log"
	"math"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TodayGamesData map[int]*GameData

func (t *TodayGamesData) getToday() string {
	return time.Now().Format("2006.01.02")
}

func (t *TodayGamesData) getCount() (count int) {
	return len(*t)
}

func (t *TodayGamesData) getMax() (max int) {
	for _, v := range *t {
		if v.level > max {
			max = v.level
		}
	}
	return max
}

func (t *TodayGamesData) getAvg() (sum float64) {
	for _, v := range *t {
		sum += float64(v.level)
	}
	if t.getCount() > 0 {
		sum /= float64(len(*t))
		return math.Round(sum*100) / 100
	}
	return 0
}

func (t *TodayGamesData) getGamesTimeDuraton() (result string) {
	if t.getCount() == 0 {
		return
	}
	dtFormat := "2006.01.02 15:04:05.000"
	var durration time.Duration
	for _, v := range *t {
		dtBeg, err := time.Parse(dtFormat, v.dtBeg)
		if err != nil {
			panic(err)
		}
		dtEnd, err := time.Parse(dtFormat, v.dtEnd)
		if err != nil {
			panic(err)
		}
		durration += dtEnd.Sub(dtBeg)
	}
	d := durration.Round(time.Millisecond)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	sec := d / time.Second
	d -= sec * time.Second
	ms := d / time.Millisecond
	if hours > 0 {
		result = fmt.Sprintf("%02v:%02v:%02v.%03v", int(hours), int(minutes), int(sec), int(ms))
	} else {
		result = fmt.Sprintf("%02v:%02v.%03v", int(minutes), int(sec), int(ms))
	}
	return
}

func (t *TodayGamesData) PlotTodayData() (gameNr, level, levelValue, percents, movesPerceent, colors list.List) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v := getApp().db.todayData[k]
		gameNr.PushBack(k)
		level.PushBack(v.level)
		result := float64(v.percent)*0.01 + float64(v.level)
		levelValue.PushBack(result)
		percents.PushBack(v.percent)
		colors.PushBack(v.BgColor())
		moves := float64(v.moves)
		totalmoves := float64(v.totalmoves)
		percentMoves := moves * 100 / totalmoves
		lvlMoves := float64(v.level) * percentMoves / 100
		movesPerceent.PushBack(lvlMoves)
	}
	return
}

func (t *TodayGamesData) getWinCountInManual() (bool, bool, int) {
	keys := make([]int, 0)
	for k := range *t {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	count := 0
	adv := getApp().preferences.manualAdv
	lastLvl := getApp().db.todayData[len(keys)].level
	ok := false
	for i := len(keys); i > 0; i-- {
		v := getApp().db.todayData[i]
		if adv == 0 || !v.manual && count < adv {
			return false, ok, count
		} else if v.manual && v.percent == 100 && v.level == lastLvl {
			count++
			ok = true
			if count == adv {
				return true, ok, count
			}
		} else if v.manual && v.percent < 100 && v.level == lastLvl {
			count = 0
			return false, ok, count
		}
		if v.level != lastLvl {
			return false, ok, count
		}
	}
	return count >= adv, ok, count
}

func (t *TodayGamesData) String() string {
	s := fmt.Sprintf("%v", t.getToday())
	if t.getCount() > 0 {
		s = fmt.Sprintf("%v #%v max:%v, avg:%v [%v]",
			t.getToday(),
			t.getCount(),
			t.getMax(),
			t.getAvg(),
			t.getGamesTimeDuraton(),
		)
	}
	return s
}

type GameData struct {
	dtBeg, dtEnd                                                                    string
	id, level, lives, percent, correct, wrong, moves, totalmoves, advance, fallback int
	manual, resetonerror                                                            bool
}

func (d *GameData) NextLevel() (int, int, string) {
	motiv := ""
	manual := getApp().preferences.manual
	adv := getApp().preferences.thresholdAdvance
	fall := getApp().preferences.thresholdFallback
	level := d.level
	lives := d.lives
	if manual {
		win, ok, count := getApp().db.todayData.getWinCountInManual()
		if !win && !ok {
			level = getApp().preferences.defaultLevel
			lives = count
		} else if !win && ok {
			motiv = "Manual game mode. Good result! One more time this level!"
			lives = count
		} else if win && ok {
			motiv = "Manual game mode. Excellent result! Level up!"
			level += 1
			lives = count
		}
	} else if d.percent >= adv {
		level += 1
		lives = getApp().preferences.thresholdFallbackSessions
		motiv = "Classic game mode. Excellent result! Level up!"
	} else if d.percent >= fall && d.percent < adv {
		motiv = "Classic game mode. Good result! One more time this level!"
	} else if d.percent < fall {
		if lives == 1 {
			motiv = "Classic game mode. Let's improve the results! Level down!"
			if level > 1 {
				level -= 1
				lives = getApp().preferences.thresholdFallbackSessions
			}
		} else if lives > 1 {
			motiv = "Classic game mode. Let's improve the results! Let's have an extra try!"
			lives -= 1
		}
	}
	return level, lives, motiv
}

func (d GameData) BgColor() (result color.Color) {
	theme := getApp().theme
	colorRegular := theme.regular
	colorCorrect := theme.correct
	colorError := theme.error
	colorWarning := theme.warning
	adv := getApp().preferences.thresholdAdvance
	fall := getApp().preferences.thresholdFallback
	if d.percent >= adv {
		result = colorRegular
	} else if d.percent >= fall && d.percent < adv {
		result = colorCorrect
	} else if d.percent < fall {
		if d.lives == 1 {
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
	dtFormat := "2006.01.02 15:04:05.000"
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
		getApp().db.todayGamesCount,
		q.level,
		q.percent,
		q.correct,
		q.wrong,
		q.moves,
		dStr)
	if getApp().preferences.resetOnFirstWrong {
		ss = fmt.Sprintf("#%v nB%v %v%% correct:%v wrong:%v moves:(%v/%v) [%v]",
			getApp().db.todayGamesCount,
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

type Db struct {
	conn            *sql.DB
	todayData       TodayGamesData
	todayGamesCount int
}

func (d *Db) Setup() {
	var err error
	d.conn, err = sql.Open("sqlite3", "games.db")
	if err != nil {
		panic(err)
	}
	var createDB string = "CREATE TABLE IF NOT EXISTS simple(id INTEGER PRIMARY KEY AUTOINCREMENT,dtBeg TEXT, dtEnd TEXT, level INTEGER, lives INTEGER, percent INTEGER, correct INTEGER, wrong NTEGER, moves INTEGER, totalmoves INTEGER, manual INTEGER, advance INTEGER, fallback INTEGER, resetonerror INTEGER)"
	cur, err := d.conn.Prepare(createDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()
}

func (d *Db) Insert(values *GameData) {
	insStr := "INSERT INTO simple(dtBeg, dtEnd, level, lives, percent, correct, wrong, moves, totalmoves, manual, advance, fallback, resetonerror) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	cur, err := d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error in DB:", insStr, values)
		panic(err)
	}
	dtBeg := values.dtBeg
	dtEnd := values.dtEnd
	level := values.level
	lives := values.lives
	percent := values.percent
	correct := values.correct
	wrong := values.wrong
	moves := values.moves
	totalmoves := values.totalmoves
	manual := values.manual
	advance := values.advance
	fallback := values.fallback
	resetonerror := values.resetonerror
	cur.Exec(dtBeg, dtEnd, level, lives, percent, correct, wrong, moves, totalmoves, manual, advance, fallback, resetonerror)
	d.todayGamesCount += 1
	d.todayData[d.todayGamesCount] = values
	log.Println("DB: Inserted:", dtBeg, dtEnd, level, lives, percent, correct, wrong, moves, totalmoves, manual, advance, fallback, resetonerror)
}

func (d *Db) ReadTodayGames() {
	d.todayData = make(map[int]*GameData)
	d.todayGamesCount = 0
	rows, err := d.conn.Query("SELECT * FROM simple")
	if err != nil {
		panic(err)
	}
	now := time.Now()
	todayBeginDt := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	dtFormat := "2006.01.02 15:04:05.000"
	for rows.Next() {
		values := &GameData{}
		err = rows.Scan(&values.id, &values.dtBeg, &values.dtEnd, &values.level, &values.lives, &values.percent, &values.correct, &values.wrong, &values.moves, &values.totalmoves, &values.manual, &values.advance, &values.fallback, &values.resetonerror)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		dt, err := time.Parse(dtFormat, values.dtBeg)
		if err != nil {
			panic(err)
		}
		if dt.After(todayBeginDt) {
			d.todayGamesCount += 1
			d.todayData[d.todayGamesCount] = values
		}
	}
}

func (d *Db) Close() {
	d.conn.Close()
	log.Println("DB Closed.")
}
