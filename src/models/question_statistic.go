package models

type QuestionStatistic struct {
	ID                 uint         `json:"id"`
	QuestionID         uint         `json:"questionId"`
	AnswerCountMap     map[uint]int `json:"answerCountMap"`
	TotalAnswer        int          `json:"totalAnswer"`
	TotalCorrectAnswer int          `json:"totalCorrectAnswer"`
	TotalWrongAnswer   int          `json:"totalWrongAnswer"`
}
