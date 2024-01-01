package models

type QuestionStatistic struct {
	ID                 uint           `json:"id"`
	QuestionID         uint           `json:"questionId"`
	AnswerCountMap     map[string]int `json:"answerCountMap"`
	ChoiceStatistics   map[string]int `json:"choiceStatistics"`
	TotalAnswer        int            `json:"totalAnswer"`
	TotalCorrectAnswer int            `json:"totalCorrectAnswer"`
	TotalWrongAnswer   int            `json:"totalWrongAnswer"`
}
