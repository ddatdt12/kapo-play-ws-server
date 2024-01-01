package utils

import "fmt"

func BuildQuestionKey(code string) string {
	return "game:" + code + ":questions"
}

func BuildGameKey(code string) string {
	return "game:" + code
}

func BuildGameStateKey(code string) string {
	return "game:" + code + ":state"
}

func BuildGameKeyMembers(code string) string {
	return "game:" + code + ":users"
}

func BuildUserAnswersKey(gameCode string, username string) string {
	return "game:" + gameCode + ":user:" + username + ":answers"
}

func BuildGameQuestionAnswers(code string, questionID uint) string {
	return "game:" + code + ":question:" + fmt.Sprint(questionID) + ":answers"
}
