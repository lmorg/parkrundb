package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var (
	fDbFileName string

	fCrawlTable  FlagCrawlTable
	fCrawlRange  FlagCrawlRange
	fCrawlAll    FlagCrawlAll
	fCrawlLatest FlagCrawlAll

	rxFlagTable *regexp.Regexp
	rxFlagRange *regexp.Regexp
	rxFlagAll   *regexp.Regexp // also used for latest
)

type FlagCrawlTable []string

func (fs *FlagCrawlTable) String() string { return fmt.Sprint(*fs) }
func (fs *FlagCrawlTable) Set(value string) error {
	if !rxFlagTable.MatchString(value) {
		return errors.New(fmt.Sprintf(`"%s" does not match format: "eventname,runnumber" (text,number)`, value))
	}
	*fs = append(*fs, value)
	return nil
}

type FlagCrawlRange []string

func (fs *FlagCrawlRange) String() string { return fmt.Sprint(*fs) }
func (fs *FlagCrawlRange) Set(value string) error {
	if !rxFlagRange.MatchString(value) {
		return errors.New(fmt.Sprintf(`"%s" does not match format: "eventname,firstrun,lastrun" (text,number,number)`, value))
	}
	*fs = append(*fs, value)
	return nil
}

type FlagCrawlAll []string

func (fs *FlagCrawlAll) String() string { return fmt.Sprint(*fs) }
func (fs *FlagCrawlAll) Set(value string) error {
	if !rxFlagAll.MatchString(value) {
		return errors.New(fmt.Sprintf(`"%s" does not match format: "eventname" (text only)`, value))
	}
	*fs = append(*fs, value)
	return nil
}

func init() {
	rxFlagTable, _ = regexp.Compile(`^([a-z]+),([0-9]+)$`)
	rxFlagRange, _ = regexp.Compile(`^([a-z]+),([0-9]+),([0-9]+)$`)
	rxFlagAll, _ = regexp.Compile(`^([a-z]+)$`)
}

func Flags() {
	flag.Usage = Usage

	flag.StringVar(&fDbFileName, "db", "parkrun.db", "")
	flag.Var(&fCrawlTable, "table", "")
	flag.Var(&fCrawlRange, "range", "")
	flag.Var(&fCrawlAll, "all", "")
	flag.Var(&fCrawlAll, "latest", "")

	flag.Parse()

	if len(fCrawlTable) == 0 && len(fCrawlRange) == 0 && len(fCrawlAll) == 0 {
		fmt.Println("No run results selected for download.")
		flag.Usage()
		os.Exit(1)
	}
}

func Usage() {
	fmt.Print(`
Usage: parkrundb [--db filename]
                 [--latest event] ...
                 [--table event,runnumber] ...
                 [--range event,firstrun,lastrun] ...
                 [--all event] ...

    --db      Sqlite3 database filename. Defaults to parkrun.db
    --latest  Returns latest table from an event. Parameters: text
    --table   Returns specific table. Parameters: text,number
    --range   Returns all tables in range inclusive. Parameters: text,number,number
    --all     Returns every table from an event. Parameters: text

Multiple table/range/all/latest flags can be used. eg downloading results across multiple events.
`)
}
