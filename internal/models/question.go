// models/question.go
package models

type Question struct {
	RawText       string
	QuestionText  string
	Answers       []string
	Difficulty    string
	CorrectAnswer int
	Explanation   string
	Explanation2  string
	Subject       string
	Topic         string
}
