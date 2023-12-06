package ws

import (
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/constants"
	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/rs/zerolog/log"
)

func router(c *Client, messageTransfer *dto.MessageTransfer) {
	log.Info().Msgf("client exist ?: : %v", c != nil)
	log.Info().Msgf("c.GameSocket.ClientSet: %v", c.GameSocket.ClientSet)

	if messageTransfer.Type == dto.MessageFirstJoin {
		c.Notify(*dto.NewMessageTransfer(dto.MessageFirstJoin, c.Game,
			map[string]interface{}{
				"stage":    c.GameSocket.GameStage,
				"status":   c.GameSocket.Status,
				"question": c.GameSocket.Question,
			},
		))
		c.GameSocket.NotifyUpdatedListPlayers()
	} else if c.IsHost {
		hostMessageHandler(c, messageTransfer)
	} else {
		playerMessageHandler(c, messageTransfer)
	}
}

func hostMessageHandler(c *Client, message *dto.MessageTransfer) {
	switch message.Type {
	case dto.MessageStartGame:
		c.GameSocket.CurrentQuestion = -1
		currentQuest := c.GameSocket.CurrentQuestion

		err := c.Hub.GameService.StartGame(c.ConnectionCtx.Ctx, c.Game.Code)
		if err != nil {
			ResponseError(c, err)
			return
		}
		question, err := c.Hub.QuestionService.GetQuestion(c.ConnectionCtx.Ctx, c.Game.ID, currentQuest+1)

		if err != nil {
			ResponseError(c, err)
			return
		}
		question.StartAt = types.NewNullableTime(time.Now().Add(constants.WaitingTimeBeforeStart * time.Second))
		c.GameSocket.CurrentQuestion = currentQuest + 1
		c.GameSocket.Question = question
		c.GameSocket.Status = models.GameStatusPlaying
		c.GameSocket.SetGameStage(models.GameStageShowQuestion)
		response := dto.MessageTransfer{
			Type: dto.MessageNewQuestion,
			Data: dto.NewQuestionRes(question),
			Meta: map[string]interface{}{
				"action": dto.MessageStartGame,
			},
		}
		c.GameSocket.NotifyAll(response)
	case dto.MessageTimeUp:
	case dto.MessageNextAction:
		if c.GameSocket.GameStage == models.GameStageShowQuestion {
			for k := range c.GameSocket.ClientSet {
				response := dto.MessageTransfer{
					Type: dto.MessageQuestionResult,
					Data: dto.NewQuestionResult(c.GameSocket.Question, k.QuestionAnswersMap[c.GameSocket.Question.ID], nil),
				}
				c.GameSocket.NotifyTo(k, response)
			}
			questionStatistic, _ := c.Hub.QuestionService.StatisticAnswersOfQuestion(c.ConnectionCtx.Ctx, c.Game.ID, c.GameSocket.Question.ID)
			c.GameSocket.SetGameStage(models.GameStageShowAnswer)
			c.GameSocket.NotifyHost(dto.MessageTransfer{
				Type: dto.MessageQuestionResult,
				Data: dto.NewQuestionResult(c.GameSocket.Question, nil, questionStatistic),
			})
		} else if c.GameSocket.GameStage == models.GameStageShowAnswer {
			c.GameSocket.SetGameStage(models.GameStageShowQuestion)
			currentQuest := c.GameSocket.CurrentQuestion
			question, err := c.Hub.QuestionService.GetQuestion(c.ConnectionCtx.Ctx, c.Game.ID, currentQuest+1)

			if err != nil {
				ResponseError(c, err)
				return
			}
			// It means game is ended
			if question == nil {
				responseGameEnded(c)
				return
			}

			question.StartAt = types.NewNullableTime(time.Now().Add(constants.WaitingTimeBeforeStart * time.Second))
			c.GameSocket.CurrentQuestion = currentQuest + 1
			c.GameSocket.Question = question
			c.GameSocket.SetGameStage(models.GameStageShowQuestion)
			response := dto.MessageTransfer{
				Type: dto.MessageNewQuestion,
				Data: dto.NewQuestionRes(question),
			}
			c.GameSocket.NotifyAll(response)
		}
	case dto.MessageLeaderboard:
		leaderBoard, err := c.Hub.LeaderboardService.GetLeaderboard(c.ConnectionCtx.Ctx, c.Game.ID)

		if err != nil {
			ResponseError(c, err)
			return
		}

		c.Send <- dto.MessageTransfer{
			Type: dto.MessageLeaderboard,
			Data: leaderBoard.Items,
		}
	}
}

func playerMessageHandler(c *Client, message *dto.MessageTransfer) {
	switch message.Type {
	case dto.MessageAnswerQuestion:
		var answerDto dto.AnswerQuestionReq
		message.Binding(&answerDto)

		log.Info().Msgf("answerDto: %v", answerDto)
		answer, err := c.Hub.QuestionService.AnwserQuestion(c.ConnectionCtx.Ctx, c.Game.ID, c.User, answerDto)
		if err != nil {
			ResponseError(c, err)
			return
		}
		c.QuestionAnswersMap[answer.QuestionID] = answer
	}
}

func ResponseError(c *Client, err error) {
	log.Error().Stack().Err(err).Msg("error")
	c.Send <- dto.MessageTransfer{
		Type: dto.Error,
		Data: err.Error(),
	}
}

func responseGameEnded(c *Client) {
	c.GameSocket.SetGameStatus(models.GameStatusEnded)
	for client := range c.GameSocket.ClientSet {
		userRank, err := c.Hub.LeaderboardService.GetUserRank(c.ConnectionCtx.Ctx, c.Game.ID, client.User.Username)
		if err != nil {
			ResponseError(client, err)
			return
		}
		c.GameSocket.NotifyTo(client, dto.MessageTransfer{
			Type: dto.MessageEndGame,
			Data: userRank,
		})
	}

	leaderBoard, err := c.Hub.LeaderboardService.GetLeaderboard(c.ConnectionCtx.Ctx, c.Game.ID)

	if err != nil {
		ResponseError(c, err)
		return
	}

	c.Send <- dto.MessageTransfer{
		Type: dto.MessageEndGame,
		Data: leaderBoard.Items,
	}

	return
}
