package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	rxTableBody *regexp.Regexp
	rxTableRow  *regexp.Regexp
	rxTableCell *regexp.Regexp
	rxStripTags *regexp.Regexp
)

func init() {
	rxTableBody, _ = regexp.Compile(`<tbody>(.*?)</tbody>`)
	rxTableRow, _ = regexp.Compile(`<tr.*?>(.*?)</tr>`)
	rxTableCell, _ = regexp.Compile(`<td.*?>(.*?)</td>`)
	rxStripTags, _ = regexp.Compile(`<.*?>`)
}

func CrawlRange(event string, runNumberFirst, runNumberLast int) {
	for i := runNumberFirst; i <= runNumberLast; i++ {
		err := Crawler(event, i)
		if err != nil && len(err.Error()) > 20 && err.Error()[:20] == "Could not find table" {
			log.Println("Assuming no more events in range")
			return
		}
	}
}

func Crawler(event string, runNumber int) (err error) {
	err = GetResults(event, runNumber)
	if err != nil {
		log.Println(err)
	}
}

func GetResults(event string, runNumber int) (err error) {
	var (
		resp *http.Response
		body []byte
	)

	const uri = "http://www.parkrun.org.uk/%s/results/weeklyresults/?runSeqNumber=%d"
	if resp, err = http.Get(fmt.Sprintf(uri, event, runNumber)); err != nil {
		return
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	err = ParseBody(&body, event, runNumber)

	return
}

func ParseBody(body *[]byte, event string, runNumber int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Panic caught: %s", r))
		}
	}()

	table := rxTableBody.FindStringSubmatch(string(*body))
	if len(table) < 2 {
		return errors.New(fmt.Sprintf("Could not find table in event %s, run# %d", event, runNumber))
	}

	rows := rxTableRow.FindAllStringSubmatch(table[1], -1)
	for i := range rows {
		var rec Record

		cells := rxTableCell.FindAllStringSubmatch(rows[i][1], -1)

		parkrunner := rxStripTags.ReplaceAllString(cells[1][1], "")

		if parkrunner == Unknown {
			log.Println(fmt.Sprintf("Skipping unknown in event %s, run# %d, row %d", event, runNumber, i))
			continue
		}

		if len(cells) < 10 {
			return errors.New(fmt.Sprintf("Failing; too few cells in event %s, run# %d, row %d", event, runNumber, i))
		}

		rec.Event = event
		rec.RunNumber = runNumber
		rec.Pos, _ = strconv.Atoi(rxStripTags.ReplaceAllString(cells[0][1], ""))
		rec.ParkRunner = parkrunner
		rec.Time = rxStripTags.ReplaceAllString(cells[2][1], "")
		rec.AgeCat = rxStripTags.ReplaceAllString(cells[3][1], "")
		rec.AgeGrade = rxStripTags.ReplaceAllString(cells[4][1], "")
		rec.Gender = rxStripTags.ReplaceAllString(cells[5][1], "")[0]
		rec.GenderPos, _ = strconv.Atoi(rxStripTags.ReplaceAllString(cells[6][1], ""))
		rec.Club = rxStripTags.ReplaceAllString(cells[7][1], "")
		rec.Note = rxStripTags.ReplaceAllString(cells[8][1], "")
		rec.TotalRuns, _ = strconv.Atoi(rxStripTags.ReplaceAllString(cells[9][1], ""))

		if err = InsertRecord(rec); err != nil {
			return
		}

	}

	return
}
