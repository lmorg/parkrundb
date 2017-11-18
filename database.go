package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	sqlCreateTable = `CREATE TABLE IF NOT EXISTS %s.results (
							id              TEXT PRIMARY KEY,
							event_code      TEXT,
							event_name      TEXT,
							run_number      INTEGER,
							date            TEXT,
							pos             INTEGER,
							parkrunner      TEXT,
							time            TEXT,
							age_cat         TEXT,
							age_grade       TEXT,
							gender          TEXT,
							gender_pos      INTEGER,
							club            TEXT,
							note            TEXT,
							total_runs      INTEGER
						);`

	sqlInsertRecord = `INSERT INTO mem.results (
							id,
							event_code,
							event_name,
							run_number,
							date,
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
						) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSyncToDisk = `INSERT INTO main.results SELECT * FROM mem.results;`
	sqlPurgeMem   = `DELETE FROM mem.results;`
)

var db *sql.DB

func OpenDB() {
	var err error

	log.Println("Opening database")

	if db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s", fDbFileName)); err != nil {
		log.Fatalln("Could not open database:", err)
	}

	if _, err = db.Exec(fmt.Sprintf(sqlCreateTable, "main")); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	if _, err = db.Exec(`ATTACH DATABASE ':memory:' AS mem;`); err != nil {
		log.Fatalln("Could not create in memory database")
	}

	if _, err = db.Exec(fmt.Sprintf(sqlCreateTable, "mem")); err != nil {
		log.Fatalln("Could not create memory table:", err)
	}
}

func InsertRecord(rec Record) (err error) {
	_, err = db.Exec(sqlInsertRecord,
		fmt.Sprintf("%s:%d:%d", rec.EventCode, rec.RunNumber, rec.Pos),
		rec.EventCode,
		rec.EventName,
		rec.RunNumber,
		rec.Date,
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

	return
}

func SyncDbToDisk() {
	if _, err := db.Exec(sqlSyncToDisk); err != nil {
		log.Println("Error syncing to disk:", err)
	}
	db.Exec(sqlPurgeMem)
	return
}

func CloseDB() {
	db.Close()
}