package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
)

func (d *Db) createAppConfTable() {
	var createGameDB string = "CREATE TABLE IF NOT EXISTS app_conf(id INTEGER PRIMARY KEY AUTOINCREMENT, fullscreen INTEGER, restperiod INTEGER, positionkey TEXT, colorkey TEXT, numberkey TEXT, ariphmetickey TEXT, lang TEXT)"
	cur, err := d.conn.Prepare(createGameDB)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	cur.Exec()
	cur.Close()
	log.Println("Created table app_conf")
}

func (d *Db) InsertAppConf() {
	if d.conn == nil {
		d.setup()
	}
	values := eui.GetUi().GetSettings()
	var empthyPreviousSettings = "DELETE from app_conf WHERE id>0"
	cur, err := d.conn.Prepare(empthyPreviousSettings)
	if err != nil {
		panic(err)
	}
	cur.Exec()
	log.Println("Deleted previous settings")

	insStr := "INSERT INTO app_conf(fullscreen, restperiod, positionkey, colorkey, numberkey, ariphmetickey, lang) VALUES(?,?,?,?,?,?,?)"
	curIns, err := d.conn.Prepare(insStr)
	if err != nil {
		log.Println("Error in DB:", insStr, values)
		panic(err)
	}
	defer curIns.Close()
	fullscreen := values.Get(eui.UiFullscreen)
	restperiod := values.Get(app.RestDuration)
	positionkey, err := values.Get(app.PositionKeypress).(ebiten.Key).MarshalText()
	if err != nil {
		panic(err)
	}
	colorkey, err := values.Get(app.ColorKeypress).(ebiten.Key).MarshalText()
	if err != nil {
		panic(err)
	}
	numberkey, err := values.Get(app.NumberKeypress).(ebiten.Key).MarshalText()
	if err != nil {
		panic(err)
	}
	ariphmetickey, err := values.Get(app.AriphmeticsKeypress).(ebiten.Key).MarshalText()
	if err != nil {
		panic(err)
	}
	lang := values.Get(app.AppLang).(string)
	curIns.Exec(fullscreen, restperiod, positionkey, colorkey, numberkey, ariphmetickey, lang)
	log.Println("DB:Inserted:", values, curIns)
}

func (d *Db) GetFromDbAppConfData() *eui.Setting {
	if d.conn == nil {
		d.setup()
	}
	rows, err := d.conn.Query("SELECT * FROM app_conf")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer rows.Close()
	conf := eui.GetUi().GetSettings()

	id := 0
	fullscreen := 0
	restperiod := 0
	positionkey := ""
	colorkey := ""
	numberkey := ""
	ariphmetickey := ""
	lang := ""

	for rows.Next() {
		err = rows.Scan(&id, &fullscreen, &restperiod, &positionkey, &colorkey, &numberkey, &ariphmetickey, &lang)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			panic(err)
		}
		fmt.Println("db read result:", id, fullscreen, restperiod, positionkey, colorkey, numberkey, ariphmetickey)
		if fullscreen == 0 {
			conf.Set(eui.UiFullscreen, false)
		} else {
			conf.Set(eui.UiFullscreen, true)
		}
		conf.Set(app.RestDuration, restperiod)

		posKey := conf.Get(app.PositionKeypress).(ebiten.Key)
		posKey.UnmarshalText([]byte(positionkey))
		conf.Set(app.PositionKeypress, posKey)

		colKey := conf.Get(app.ColorKeypress).(ebiten.Key)
		colKey.UnmarshalText([]byte(colorkey))
		conf.Set(app.ColorKeypress, colKey)

		numKey := conf.Get(app.NumberKeypress).(ebiten.Key)
		numKey.UnmarshalText([]byte(numberkey))
		conf.Set(app.NumberKeypress, numKey)

		ariKey := conf.Get(app.AriphmeticsKeypress).(ebiten.Key)
		ariKey.UnmarshalText([]byte(ariphmetickey))
		conf.Set(app.AriphmeticsKeypress, ariKey)
		conf.Set(app.AppLang, lang)
	}
	if len(*conf) > 0 {
		log.Println("Done Read AppConf table from DB", conf, len(*conf), "items")
		return conf
	}
	log.Println("Read app_conf nil from DB")
	return nil
}
