package main

import (
	"log"
)

const Event = "" // event name as found in the URL. eg "hampsteadheath"
const Filename = "parkrun.db"
const Version = "0.1 ALPHA"

func main() {
	log.Println("Version", Version)

	OpenDB(Filename)

	CrawlRange(Event, 1, 130)
	//Crawler(Event, 1)

	if err := SyncDbToDisk(Filename); err != nil {
		log.Println(err)
	}

	CloseDB()
}
