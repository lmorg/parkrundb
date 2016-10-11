# parkrundb
Creates or updates an sqlite3 database from the Parkrun Results page

## Flags:

    Usage: parkrundb [--db filename]
                     [--table event,runnumber] ...
                     [--range event,firstrun,lastrun] ...
                     [--all event] ...

        --db      Sqlite3 database filename. Defaults to parkrun.db
        --table   Returns specific table. Parameters: text,number
        --range   Returns all tables in range inclusive. Parameters: text,number,number
        --all     Returns every table from an event. Parameters: text

    Multiple table/range/all flags can be used. eg downloading results across multiple events.

## Prerequisites:

If you haven't already, you will need the Go (golang) toolchain installed on your machine to compile this source code: https://golang.org/

## Install instructions:

### Linux, OS X, FreeBSD:

	go get github.com/mattn/go-sqlite3
	go install github.com/mattn/go-sqlite3
	go install github.com/lmorg/parkrundb


### Windows install notes:
In addition to the Go language, you will need gcc installed to run `go install` against sqlite3:
https://sourceforge.net/projects/mingw-w64/?source=typ_redirect

Also you will need git installed (if it isn't already):
https://git-scm.com/download/win

Then run:

	set PATH=%PATH%;c:\Program Files\mingw-w64\x86_64-6.2.0-posix-seh-rt_v5-rev1\mingw64\bin

(where the above path is the install destination of mingw-w64)

	go get github.com/mattn/go-sqlite3
	go install github.com/mattn/go-sqlite3
	go install github.com/lmorg/parkrundb

## Recompiling changes to _parkrundb_:

Simply run:

	go install github.com/lmorg/parkrundb
