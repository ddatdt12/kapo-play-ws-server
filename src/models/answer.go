package models

type Answer struct {
	Choices    []any `json:"choices"`
	QuestionID uint  `json:"questionId"`
	IsCorrect  bool  `json:"isCorrect"`
	AnswerTime int64 `json:"answerTime"`
	Point      int64 `json:"point"`
}
