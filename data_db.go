package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	conn            *sql.DB
	todayData       TodayGamesData
	todayGamesCount int
	scoresData      ScoresData
}

func (d *Db) Setup() {
	var err error
	d.conn, err = sql.Open("sqlite3", "games.db")
	if err != nil {
		panic(err)
	}
	var createGameDB string = "CREATE TABLE IF NOT EXISTS simple(id INTEGER PRIMARY KEY AUTOINCREMENT,dtBeg TEXT, dtEnd TEXT, level INTEGER, lives INTEGER, percent INTEGER, correct INTEGER, wrong NTEGER, moves INTEGER, totalmoves INTEGER, manual INTEGER, advance INTEGER, fallback INTEGER, resetonerror INTEGER)"
	cur, err := d.conn.Prepare(createGameDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()

	var createSettingsDB string = "CREATE TABLE IF NOT EXISTS settings(id INTEGER PRIMARY KEY AUTOINCREMENT,timetonextcell REAL, timeshowcell REAL, rr REAL, level INTEGER, manualadv INTEGER,manual INTEGER, thresholdadv INTEGER, thresholdfall INTEGER, threshholssessions INTEGER, trials INTEGER, factor INTEGER, exponent INTEGER, feedbackmove INTEGER, usecenter INTEGER, resetonwrong INTEGER, fullscreen INTEGER, pauserest INTEGER, grid INTEGER)"
	cur, err = d.conn.Prepare(createSettingsDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()
}

func (d *Db) InsertSettings(values *Setting) {
	var empthyPreviousSettings = "DELETE from settings WHERE id>0"
	cur, err := d.conn.Prepare(empthyPreviousSettings)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	log.Println("Deleted previous settings")

	insStr := "INSERT INTO settings(timetonextcell, timeshowcell, rr, level, manualadv, manual, thresholdadv, thresholdfall, threshholssessions, trials, factor, exponent, feedbackmove, usecenter, resetonwrong, fullscreen, pauserest, grid) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	cur, err = d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error in DB:", insStr, values)
		panic(err)
	}
	cur.Exec(values.TimeToNextCell, values.TimeShowCell, values.RR, values.DefaultLevel, values.ManualAdv, values.Manual, values.ThresholdAdvance, values.ThresholdFallback, values.ThresholdFallbackSessions, values.Trials, values.TrialsFactor, values.TrialsExponent, values.FeedbackOnUserMove, values.Usecentercell, values.ResetOnFirstWrong, values.FullScreen, values.PauseRest, values.GridSize)
	log.Println("DB: Inserted settings.")
}

func (d *Db) ReadSettings(values *Setting) {
	if values.DefaultLevel == 0 {
		values.Reset()
	}
	rows, err := d.conn.Query("SELECT * FROM settings")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		id := 0
		err = rows.Scan(&id, &values.TimeToNextCell, &values.TimeShowCell, &values.RR, &values.DefaultLevel, &values.ManualAdv, &values.Manual, &values.ThresholdAdvance, &values.ThresholdFallback, &values.ThresholdFallbackSessions, &values.Trials, &values.TrialsFactor, &values.TrialsExponent, &values.FeedbackOnUserMove, &values.Usecentercell, &values.ResetOnFirstWrong, &values.FullScreen, &values.PauseRest, &values.GridSize)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
	}
}

func (d *Db) InsertGame(values *GameData) {
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
	dtFormat := "2006-01-02 15:04:05.000"
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

func (d *Db) ReadAllGamesScore() (*ScoreData, string) {
	qry := "SELECT count(level) games, max(level) max, round(avg(level),2) average FROM simple;"
	rows, err := d.conn.Query(qry)
	if err != nil {
		panic(err)
	}
	values := &ScoreData{}
	resultStr := ""
	for rows.Next() {
		err = rows.Scan(&values.games, &values.max, &values.avg)
		if values.games == 0 {
			break
		}
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		resultStr = fmt.Sprintf("Games: %v Max:%v Average:%v", values.games, values.max, values.avg)
	}
	return values, resultStr
}

func (d *Db) ReadAllGamesForScoresByDays() {
	qry := "SELECT count() games,max(level)max,round( avg(level),2)average, strftime('%Y-%m-%d',datetime(dtBeg)) day FROM simple GROUP BY day;"
	d.scoresData = make(ScoresData)
	rows, err := d.conn.Query(qry)
	if err != nil {
		panic(err)
	}
	dtFormat := "2006-01-02"
	i := 1
	for rows.Next() {
		values := &ScoreData{}
		var dStr string
		err = rows.Scan(&values.games, &values.max, &values.avg, &dStr)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		dt, err := time.Parse(dtFormat, dStr)
		if err != nil {
			panic(err)
		}
		values.dt = dt
		d.scoresData[i] = values
		i++
	}
}

func (d *Db) Close() {
	d.conn.Close()
	log.Println("DB Closed.")
}
