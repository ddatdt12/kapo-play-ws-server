package models

import (
	"time"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeFillInBlank    QuestionType = "fill_in_blank"
	QuestionTypeTypeAnswer     QuestionType = "type_answer"
	QuestionOpenEnded          QuestionType = "open_ended"
	QuestionTypePoll           QuestionType = "poll"
	QuestionWordCloud          QuestionType = "word_cloud"
)

type Question struct {
	ID         uint
	Content    string
	Type       QuestionType
	TemplateID uint
	LimitTime  uint
	Points     uint
	Choices    []QuestionChoice
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type QuestionChoice struct {
	Content   string
	IsCorrect bool
}

func (s QuestionChoice) ValidTypes() []QuestionType {
	return []QuestionType{QuestionTypeMultipleChoice, QuestionTypeFillInBlank, QuestionTypePoll, QuestionWordCloud}
}

func (s QuestionType) IsValid() bool {
	switch s {
	case QuestionTypeMultipleChoice, QuestionTypeFillInBlank, QuestionOpenEnded, QuestionTypePoll, QuestionWordCloud:
		return true
	}

	return false
}
