package main

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const uriTableRunNumber = "http://www.parkrun.org.uk/%s/results/weeklyresults/?runSeqNumber=%d"
const uriTableLatest = "http://www.parkrun.org.uk/%s/results/latestresults/"

const errNoTable = "Could not find table"

const latest = 0

var (
	rxPageTitle *regexp.Regexp
	rxTableBody *regexp.Regexp
	rxTableRow  *regexp.Regexp
	rxTableCell *regexp.Regexp
	rxStripTags *regexp.Regexp
)

func init() {
	rxPageTitle, _ = regexp.Compile(`<h2>(.*?) parkrun #\s+([0-9]+) -\s+([0-9]{2})/([0-9]{2})/([0-9]{4})</h2>`)
	rxTableBody, _ = regexp.Compile(`<tbody>(.*?)</tbody>`)
	rxTableRow, _ = regexp.Compile(`<tr.*?>(.*?)</tr>`)
	rxTableCell, _ = regexp.Compile(`<td.*?>(.*?)</td>`)
	rxStripTags, _ = regexp.Compile(`<.*?>`)
}

func CrawlTable(event string, runNumber int) (err error) {
	log.Printf("Processing event %s #%d", event, runNumber)

	if err = GetResults(event, runNumber); err != nil {
		log.Println(err)
	}

	return
}

func CrawlRange(event string, firstRun, lastRun int) {
	for i := firstRun; i <= lastRun; i++ {
		err := CrawlTable(event, i)
		if checkErrNoTable(err) {
			log.Printf("Assuming no more runs in range: event %s, firstRun %d, lastRun %d", event, firstRun, lastRun)
			return
		}
	}
}

func CrawlAll(event string) {
	for i := 1; !checkErrNoTable(CrawlTable(event, i)); i++ {
		// No need for code here because the function call is also used as the loop exception:
		// !checkErrNoTable(CrawlTable(...))
	}
	log.Printf("Assuming no more runs in this event: %s", event)
}

func checkErrNoTable(err error) bool {
	if err != nil && !strings.HasSuffix(err.Error(), errNoTable) {
		return true
	}
	return false
}

func GetResults(event string, runNumber int) (err error) {
	var (
		uri  string
		resp *http.Response
		body []byte
	)

	if runNumber == latest {
		uri = fmt.Sprintf(uriTableLatest, event)
	} else {
		uri = fmt.Sprintf(uriTableRunNumber, event, runNumber)
	}

	if resp, err = http.Get(uri); err != nil {
		return
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	err = ParseBody(string(body), event, runNumber)

	return
}

func ParseBody(body string, event string, runNumber int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Panic caught: %s", r))
		}
	}()

	BeginTransaction()
	defer CommitTransaction()

	table := rxTableBody.FindStringSubmatch(body)
	if len(table) < 2 {
		return errors.New(fmt.Sprintf("%s %s, run# %d", errNoTable, event, runNumber))
	}

	title := rxPageTitle.FindStringSubmatch(body)
	eventName := title[1]
	runDate := title[5] + "-" + title[4] + "-" + title[3] // inside loop so concatenation is faster than fmt.Sprinf()
	if runNumber == latest {
		runNumber, _ = strconv.Atoi(title[2])
	}

	rows := rxTableRow.FindAllStringSubmatch(table[1], -1)
	for i := range rows {
		var rec Record

		cells := rxTableCell.FindAllStringSubmatch(rows[i][1], -1)

		parkrunner := cleanseRecord(cells[1][1])

		if parkrunner == Unknown {

			rec.EventCode = event
			rec.EventName = eventName
			rec.RunNumber = runNumber
			rec.Date = runDate
			rec.Pos, _ = strconv.Atoi(cleanseRecord(cells[0][1]))
			rec.ParkRunner = parkrunner
			rec.Gender = "U"
			if err = InsertRecord(rec); err != nil {
				return
			}
			continue
		}

		if len(cells) < 10 {
			return errors.New(fmt.Sprintf("Failing; too few cells in event %s, run# %d, row %d", event, runNumber, i))
		}

		rec.EventCode = event
		rec.EventName = eventName
		rec.RunNumber = runNumber
		rec.Date = runDate
		rec.Pos, _ = strconv.Atoi(cleanseRecord(cells[0][1]))
		rec.ParkRunner = parkrunner
		rec.Time = cleanseRecord(cells[2][1])
		rec.AgeCat = cleanseRecord(cells[3][1])
		rec.AgeGrade = cleanseRecord(cells[4][1])
		rec.Gender = cleanseRecord(cells[5][1])
		rec.GenderPos, _ = strconv.Atoi(cleanseRecord(cells[6][1]))
		rec.Club = cleanseRecord(cells[7][1])
		rec.Note = cleanseRecord(cells[8][1])
		rec.TotalRuns, _ = strconv.Atoi(cleanseRecord(cells[9][1]))
		if err = InsertRecord(rec); err != nil {
			return
		}

	}

	return
}

func cleanseRecord(s string) string {
	s = rxStripTags.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	return s
}
