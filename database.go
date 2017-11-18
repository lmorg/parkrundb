package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	sqlCreateTable = `CREATE TABLE IF NOT EXISTS results (
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

	sqlInsertRecord = `INSERT INTO results (
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
)

var (
	db *sql.DB
	tx *sql.Tx
)

func OpenDB() {
	var err error

	log.Println("Opening database")

	if db, err = sql.Open("sqlite3", "file:"+fDbFileName); err != nil {
		log.Fatalln("Could not open database:", err)
	}

	if _, err = db.Exec(sqlCreateTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	/*if _, err = db.Exec(`ATTACH DATABASE ':memory:' AS mem;`); err != nil {
		log.Fatalln("Could not create in memory database")
	}

	if _, err = db.Exec(fmt.Sprintf(sqlCreateTable, "mem")); err != nil {
		log.Fatalln("Could not create memory table:", err)
	}*/
}

func InsertRecord(rec Record) (err error) {
	_, err = tx.Exec(sqlInsertRecord,
		fmt.Sprintf("%s:%d:%d", rec.EventCode, rec.RunNumber, rec.Pos), // unique key
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

func BeginTransaction() {
	var err error
	if tx, err = db.Begin(); err != nil {
		log.Fatalln("Could not open transaction:", err)
	}
}

func CommitTransaction() {
	if err := tx.Commit(); err != nil {
		log.Println("Error commiting transaction:", err)
	}
}

func CloseDB() {
	db.Close()
}
