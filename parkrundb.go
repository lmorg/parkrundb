package main

import (
	"log"
	"strconv"
)

const Version = "1.4"

func main() {
	log.Println("Version", Version)

	Flags()

	OpenDB()

	for _, flag := range fCrawlTable {
		CrawlTable(ParsefTable(flag))
	}

	for _, flag := range fCrawlRange {
		CrawlRange(ParsefRange(flag))
	}

	for _, flag := range fCrawlAll {
		CrawlAll(flag)
	}

	for _, flag := range fCrawlLatest {
		CrawlTable(flag, latest)
	}

	CloseDB()
}

func ParsefTable(flag string) (event string, runNumber int) {
	params := rxFlagTable.FindStringSubmatch(flag)
	event = params[1]
	runNumber, _ = strconv.Atoi(params[2])
	return
}

func ParsefRange(flag string) (event string, firstRun, lastRun int) {
	params := rxFlagRange.FindStringSubmatch(flag)
	event = params[1]
	firstRun, _ = strconv.Atoi(params[2])
	lastRun, _ = strconv.Atoi(params[3])
	return
}
