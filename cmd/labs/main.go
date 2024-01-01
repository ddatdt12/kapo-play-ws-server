package main

import (
	"flag"

	"github.com/ddatdt12/kapo-play-ws-server/clients/answer"
	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/ddatdt12/kapo-play-ws-server/internal/logger"
)

func main() {
	flag.Parse()
	logger.InitLogger()
	configs.InitEnvConfigs()

	answerClient := answer.NewAnswerClient()
	// answerClient.CreateAnswer(&answer.CreateAnswer{
	// 	QuestionID: 1,
	// 	Values:     []any{1, 2, 3},
	// 	GameID:     1,
	// 	IsCorrect:  true,
	// 	Points:     10,
	// 	AnswerTime: 5,
	// 	AnsweredAt: time.Now(),
	// 	Username:   "test",
	// })
	answerClient.GetTemplates()
}
