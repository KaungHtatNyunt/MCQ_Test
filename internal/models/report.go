// models/report.go
package models

import "time"

type Report struct {
	NumQuestions    int
	Topic           string
	CorrectAnswers  int
	ScorePercentage float64
	TimeTaken       time.Duration
	AvgTimePerQ     float64
	StartTime       time.Time
	EndTime         time.Time
}

type TestSession struct {
	Questions    []Question
	CurrentIndex int
	NumQuestions int
	CorrectCount int
	StartTime    time.Time
	EndTime      time.Time
	Answers      []int // Store student answers
}
