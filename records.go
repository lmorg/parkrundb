package main

const Unknown = "Unknown"

type Record struct {
	Event      string
	RunNumber  int
	Pos        int
	ParkRunner string
	Time       string //time.Time
	AgeCat     string
	AgeGrade   string
	Gender     byte
	GenderPos  int
	Club       string
	Note       string
	TotalRuns  int
}
