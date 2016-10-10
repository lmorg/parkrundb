package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var (
	//dbDisk *sql.DB
	//dbMem  *sql.DB
	db    *sql.DB
	mutex = &sync.Mutex{}
)

const (
	sqlCreateTable = `CREATE TABLE results (
							id			integer PRIMARY KEY,
							event		string,
							run_number	integer,
							pos			integer,
							parkrunner	string,
							time		string,
							age_cat		string,
							age_grade	string,
							gender		char,
							gender_pos	integer,
							club		string,
							note		string,
							total_runs	integer
						);`

	sqlCreateTableM = `CREATE TABLE mem.results (
							id			integer PRIMARY KEY,
							event		string,
							run_number	integer,
							pos			integer,
							parkrunner	string,
							time		string,
							age_cat		string,
							age_grade	string,
							gender		char,
							gender_pos	integer,
							club		string,
							note		string,
							total_runs	integer
						);`

	sqlInsertRecord = `INSERT INTO mem.results (
							event,
							run_number,
							pos,
							parkrunner,
							time,
							age_cat,
							age_grade,
							gender,
							gender_pos,
							club,
							note,
							total_runs
						) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
)

func OpenDB(filename string) {
	var err error

	log.Println("Opening database")

	//db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", filename))
	db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s", filename))
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	_, err = db.Exec(sqlCreateTable)
	if err != nil {
		log.Println("Could not create table:", err)
	}

	_, err = db.Exec(`ATTACH DATABASE ':memory:' AS mem;`)
	if err != nil {
		log.Fatalln("Could not create in memory database")
	}

	_, err = db.Exec(sqlCreateTableM)
	if err != nil {
		log.Fatalln("Could not create memory table:", err)
	}
}

func InsertRecord(rec Record) (err error) {
	mutex.Lock()
	_, err = db.Exec(sqlInsertRecord,
		rec.Event,
		rec.RunNumber,
		rec.Pos,
		rec.ParkRunner,
		rec.Time,
		rec.AgeCat,
		rec.AgeGrade,
		rec.Gender,
		rec.GenderPos,
		rec.Club,
		rec.Note,
		rec.TotalRuns,
	)
	mutex.Unlock()

	if err == nil {
		log.Println(rec.Event, rec.RunNumber, rec.Pos, rec.ParkRunner)
	}

	return
}

func SyncDbToDisk(filename string) (err error) {
	log.Println("Syncing memory to", filename)
	_, err = db.Exec(`INSERT INTO main.results SELECT * FROM mem.results;`)
	return
}

/*
func Query(sql string) (rows *sql.Rows, err error) {
	rows, err = db.Query(sql)

	return
}
*/

func CloseDB() {
	db.Close()
}
