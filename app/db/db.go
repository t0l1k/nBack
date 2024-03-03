package db

import (
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/quasilyte/gdata"

	"database/sql"
	"log"
)

type Db struct {
	conn *sql.DB
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

func (d *Db) setup() {
	d.openDb()
	d.createAppConfTable()
}

func (d *Db) openDb() {
	var err error
	filename := d.getDbPath()
	d.conn, err = sql.Open("sqlite3", filename)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

// Получить рабочий каталог текущей ОС для БД
func (*Db) getDbPath() string {
	m, err := gdata.Open(gdata.Config{AppName: "nback"})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	s := []byte("Working memory train")
	if err := m.SaveItem("nback.txt", s); err != nil {
		log.Println(err)
		panic(err)
	}
	filename := fmt.Sprintf("%v/nback.db", strings.TrimSuffix(m.ItemPath("nback.txt"), "/nback.txt"))
	log.Println("to db path", strings.TrimSuffix(m.ItemPath("nback.txt"), "/nback.txt"), filename)
	return filename
}
