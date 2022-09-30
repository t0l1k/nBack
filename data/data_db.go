package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/t0l1k/nBack/ui"
)

type Db struct {
	conn            *sql.DB
	TodayData       TodayGamesData
	TodayGamesCount int
	ScoresData      ScoresData
}

var dbInstance *Db = nil

func init() {
	dbInstance = GetDb()
}

func GetDb() (db *Db) {
	if dbInstance == nil {
		db = &Db{}
	} else {
		db = dbInstance
	}
	return db
}

func (d *Db) Setup() {
	d.createGamesTable()
	d.createSettingsTable()
}

func (d *Db) createGamesTable() {
	var err error
	d.conn, err = sql.Open("sqlite3", "games.db")
	if err != nil {
		panic(err)
	}
	var createGameDB string = "CREATE TABLE IF NOT EXISTS simple(id INTEGER PRIMARY KEY AUTOINCREMENT,gametype TEXT, dtBeg TEXT, dtEnd TEXT, level INTEGER, lives INTEGER, percent INTEGER, correct INTEGER, wrong INTEGER, missed INTEGER, moves INTEGER, totalmoves INTEGER, manual INTEGER, advance INTEGER, fallback INTEGER, resetonerror INTEGER)"
	cur, err := d.conn.Prepare(createGameDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	cur.Close()
	log.Println("Created table for games.")
}

func (d *Db) createSettingsTable() {
	var createSettingsDB string = "CREATE TABLE IF NOT EXISTS settings(id INTEGER PRIMARY KEY AUTOINCREMENT,timetonextcell REAL, timeshowcell REAL, rr REAL, level INTEGER, manualadv INTEGER,manual INTEGER, thresholdadv INTEGER, thresholdfall INTEGER, threshholssessions INTEGER, trials INTEGER, factor INTEGER, exponent INTEGER, feedbackmove INTEGER, usecenter INTEGER, resetonwrong INTEGER, fullscreen INTEGER, pauserest INTEGER, grid INTEGER, showgrid INTEGER,showcrosshair INTEGER, gametype TEXT, symcount INTEGER, lang TEXT)"
	cur, err := d.conn.Prepare(createSettingsDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	cur.Close()
	log.Println("Created table for settings.")
}

func (d *Db) InsertSettings(values *ui.Preferences) {
	if d.conn == nil {
		d.Setup()
	}
	var empthyPreviousSettings = "DELETE from settings WHERE id>0"
	cur, err := d.conn.Prepare(empthyPreviousSettings)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	log.Println("Deleted previous settings")

	insStr := "INSERT INTO settings(timetonextcell, timeshowcell, rr, level, manualadv, manual, thresholdadv, thresholdfall, threshholssessions, trials, factor, exponent, feedbackmove, usecenter, resetonwrong, fullscreen, pauserest, grid, showgrid, showcrosshair, gametype, symcount, lang) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	cur, err = d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error insert settings in DB:", insStr, values)
		d.dropSettingsTable()
		panic(err)
	}
	cur.Exec(
		values.Get("time to next cell").(float64),
		values.Get("time to show cell").(float64),
		values.Get("random repition").(float64),
		values.Get("default level").(int),
		values.Get("manual advance").(int),
		values.Get("manual mode").(bool),
		values.Get("threshold advance").(int),
		values.Get("threshold fallback").(int),
		values.Get("threshold fallback sessions").(int),
		values.Get("trials").(int),
		values.Get("trials factor").(int),
		values.Get("trials exponent").(int),
		values.Get("feedback on user move").(bool),
		values.Get("use center cell").(bool),
		values.Get("reset on first wrong").(bool),
		values.Get("fullscreen").(bool),
		values.Get("pause to rest").(int),
		values.Get("grid size").(int),
		values.Get("show grid").(bool),
		values.Get("show crosshair").(bool),
		values.Get("game type").(string),
		values.Get("symbols count").(int),
		values.Get("lang").(string),
	)
	log.Println("DB: Inserted settings.", values)
}

func (d *Db) dropSettingsTable() {
	log.Println("Old settings found in db!")
	var createSettingsDB string = "DROP TABLE settings;"
	cur, err := d.conn.Prepare(createSettingsDB)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	cur.Close()
	log.Println("Drop table settings done.")
}

func (d *Db) ReadSettings() (values *ui.Preferences) {
	if d.conn == nil {
		d.Setup()
	}
	rows, err := d.conn.Query("SELECT * FROM settings")
	if err != nil {
		panic(err)
	}
	v := ui.NewPreferences()
	for rows.Next() {
		id := 0
		tmnc := 0.1
		tmsnc := 0.1
		rr := 0.1
		dfl := 0
		madv := 0
		mm := false
		tda := 0
		tdaa := 0
		tfadvs := 0
		tr := 0
		trf := 0
		tre := 0
		fb := false
		ucc := false
		rfw := false
		fsc := false
		pr := 0
		gs := 0
		shgz := false
		shch := false
		gt := ""
		sc := 0
		lg := ""
		err = rows.Scan(
			&id,
			&tmnc,
			&tmsnc,
			&rr,
			&dfl,
			&madv,
			&mm,
			&tda,
			&tdaa,
			&tfadvs,
			&tr,
			&trf,
			&tre,
			&fb,
			&ucc,
			&rfw,
			&fsc,
			&pr,
			&gs,
			&shgz,
			&shch,
			&gt,
			&sc,
			&lg,
		)
		if err != nil && err != sql.ErrNoRows {
			d.dropSettingsTable()
			return nil
		}
		v.Set("time to next cell", tmnc)
		v.Set("time to show cell", tmsnc)
		v.Set("random repition", rr)
		v.Set("default level", dfl)
		v.Set("manual advance", madv)
		v.Set("manual mode", mm)
		v.Set("threshold advance", tda)
		v.Set("threshold fallback", tdaa)
		v.Set("threshold fallback sessions", tfadvs)
		v.Set("trials", tr)
		v.Set("trials factor", trf)
		v.Set("trials exponent", tre)
		v.Set("feedback on user move", fb)
		v.Set("use center cell", ucc)
		v.Set("reset on first wrong", rfw)
		v.Set("fullscreen", fsc)
		v.Set("pause to rest", pr)
		v.Set("grid size", gs)
		v.Set("show grid", shgz)
		v.Set("show crosshair", shch)
		v.Set("game type", gt)
		v.Set("symbols count", sc)
		v.Set("lang", lg)
	}
	log.Println("Read settings from db")
	if len(v) > 0 {
		return &v
	}
	return nil
}

func (d *Db) InsertGame(values *GameData) {
	if d.conn == nil {
		d.Setup()
	}
	insStr := "INSERT INTO simple(gametype, dtBeg, dtEnd, level, lives, percent, correct, wrong, missed, moves, totalmoves, manual, advance, fallback, resetonerror) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	cur, err := d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error in DB:", insStr, values)
		panic(err)
	}
	gametype := values.GameType
	dtBeg := values.DtBeg
	dtEnd := values.DtEnd
	level := values.Level
	lives := values.Lives
	percent := values.Percent
	correct := values.Correct
	wrong := values.Wrong
	missed := values.Missed
	moves := values.Moves
	totalmoves := values.Totalmoves
	manual := values.Manual
	advance := values.Advance
	fallback := values.Fallback
	resetonerror := values.Resetonerror
	cur.Exec(gametype, dtBeg, dtEnd, level, lives, percent, correct, wrong, missed, moves, totalmoves, manual, advance, fallback, resetonerror)
	d.TodayGamesCount += 1
	d.TodayData[d.TodayGamesCount] = values
	log.Println("DB: Inserted:", gametype, dtBeg, dtEnd, level, lives, percent, correct, wrong, missed, moves, totalmoves, manual, advance, fallback, resetonerror)
}

func (d *Db) ReadTodayGames() {
	d.TodayData = make(map[int]*GameData)
	d.TodayGamesCount = 0
	if d.conn == nil {
		return
	}
	rows, err := d.conn.Query("SELECT * FROM simple")
	if err != nil {
		panic(err)
	}
	now := time.Now()
	todayBeginDt := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	dtFormat := "2006-01-02 15:04:05.000"
	for rows.Next() {
		values := &GameData{}
		err = rows.Scan(&values.Id, &values.GameType, &values.DtBeg, &values.DtEnd, &values.Level, &values.Lives, &values.Percent, &values.Correct, &values.Wrong, &values.Missed, &values.Moves, &values.Totalmoves, &values.Manual, &values.Advance, &values.Fallback, &values.Resetonerror)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		dt, err := time.Parse(dtFormat, values.DtBeg)
		if err != nil {
			panic(err)
		}
		if dt.After(todayBeginDt) {
			d.TodayGamesCount += 1
			d.TodayData[d.TodayGamesCount] = values
		}
	}
}

func (d *Db) ReadAllGamesScore() (*ScoreData, string) {
	values := &ScoreData{}
	resultStr := ui.GetLocale().Get("scrResultNil")
	if d.conn == nil {
		return values, resultStr
	}
	qry := "SELECT count(level) games, max(level) max, round(avg(level),2) average FROM simple;"
	rows, err := d.conn.Query(qry)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&values.Games, &values.Max, &values.Avg)
		if values.Games == 0 {
			break
		}
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		resultStr = fmt.Sprintf("%v: %v %v:%v %v:%v",
			ui.GetLocale().Get("scrResultTtl"),
			values.Games,
			ui.GetLocale().Get("wordMax"),
			values.Max,
			ui.GetLocale().Get("wordAvg"),
			values.Avg)
	}
	return values, resultStr
}

func (d *Db) ReadAllGamesForScoresByDays() {
	if d.conn == nil {
		return
	}
	qry := "SELECT count() games,max(level)max,round( avg(level),2)average, strftime('%Y-%m-%d',datetime(dtBeg)) day FROM simple GROUP BY day;"
	d.ScoresData = make(ScoresData)
	rows, err := d.conn.Query(qry)
	if err != nil {
		panic(err)
	}
	dtFormat := "2006-01-02"
	i := 1
	for rows.Next() {
		values := &ScoreData{}
		var dStr string
		err = rows.Scan(&values.Games, &values.Max, &values.Avg, &dStr)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		dt, err := time.Parse(dtFormat, dStr)
		if err != nil {
			panic(err)
		}
		values.Dt = dt
		d.ScoresData[i] = values
		i++
	}
}

func (d *Db) Close() {
	d.conn.Close()
	log.Println("DB Closed.")
}
