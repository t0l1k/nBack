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

func (t *TodayGamesData) getTimeDuraton() (result string) {
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
	mSec := durration.Milliseconds() / 1e3
	sec := durration.Seconds()
	minutes := int(sec / 60)
	seconds := int(sec) % 60
	result = fmt.Sprintf("%02v:%02v.%03v", minutes, seconds, int(mSec))
	if sec > 3600 {
		result = fmt.Sprintf("%02v:%02v:%02v.%03v", sec/3600, minutes, seconds, int(mSec))
	}
	return
}

func (t *TodayGamesData) PlotTodayData() (gameNr, level, levelValue, percents, colors list.List) {
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
	}
	return
}

func (t *TodayGamesData) String() string {
	s := fmt.Sprintf("%v", t.getToday())
	if t.getCount() > 0 {
		s = fmt.Sprintf("%v #%v max:%v, avg:%v [%v]",
			t.getToday(),
			t.getCount(),
			t.getMax(),
			t.getAvg(),
			t.getTimeDuraton(),
		)
	}
	return s
}

// dtBeg, dtEnd, level, lives, percent, correct, wrong, moves, totalmoves, manual, advance, fallback, resetonerror
type GameData struct {
	dtBeg, dtEnd                                                                    string
	id, level, lives, percent, correct, wrong, moves, totalmoves, advance, fallback int
	manual, resetonerror                                                            bool
}

func (d *GameData) NextLevel() (int, int, string) {
	motiv := ""
	adv := 80
	fall := 50
	level := d.level
	lives := d.lives
	if d.percent >= adv {
		level += 1
		lives = 3
		motiv = "Excellent result! Level up!"
	} else if d.percent >= fall && d.percent < adv {
		motiv = "Good result! One more time this level!"
	} else if d.percent < fall {
		if lives == 1 {
			motiv = "Let's improve the results! Level down!"
			if level > 1 {
				level -= 1
				lives = 3
			}
		} else if lives > 1 {
			motiv = "Let's improve the results! Let's have an extra try!"
			lives -= 1
		}
	}
	return level, lives, motiv
}

func (d GameData) BgColor() (result color.RGBA) {
	colorRegular := color.RGBA{0, 0, 128, 255}
	colorCorrect := color.RGBA{0, 128, 0, 255}
	colorError := color.RGBA{255, 0, 0, 255}
	colorWarning := color.RGBA{255, 128, 0, 255}
	adv := 80
	fall := 50
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
	i := 1
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
		i += 1
	}
}

func (d *Db) Close() {
	d.conn.Close()
	log.Println("DB Closed.")
}
