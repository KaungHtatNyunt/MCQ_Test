// utils/excel.go
package utils

import (
	"github.com/xuri/excelize/v2"
	"mcq-test-system/internal/models"
	"strconv"
)

func LoadQuestionsFromExcel(filename string) ([]models.Question, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var questions []models.Question
	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}

		correctAnswer, _ := strconv.Atoi(row[7])
		question := models.Question{
			RawText:       row[0],
			QuestionText:  row[1],
			Answers:       []string{row[2], row[3], row[4], row[5]},
			Difficulty:    row[6],
			CorrectAnswer: correctAnswer,
			Explanation:   row[8],
			Explanation2:  row[9],
			Subject:       row[10],
			Topic:         row[11],
		}
		questions = append(questions, question)
	}
	return questions, nil
}
