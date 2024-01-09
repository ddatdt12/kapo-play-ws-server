package ws

import (
	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func router(c *Client, messageTransfer *dto.MessageTransfer) {
	log.Info().Msgf("router client User : %v", c.User)

	if messageTransfer.Type == dto.MessageFirstJoin {
		log.Info().Msgf("game state: %v", *c.GameSocket.GameState)

		if c.GameSocket.GameState.Status == models.GameStatusEnded {
			if c.IsHost {
				c.NotifyLeaderBoard(dto.MessageEndGame)
			} else {
				c.NotifyUserRank(dto.MessageEndGame)
			}
		}

		response := map[string]interface{}{
			"gameState":   c.GameSocket.GameState,
			"currentUser": c.User,
		}

		if c.IsHost {
			response["question"] = c.GameSocket.GameState.Question
		}
		if c.GameSocket.GameState.Question != nil {
			response["questionStatus"] = c.GameSocket.GameState.Question.Status
		}

		log.Debug().Interface("response first join: ", response).Msg("response first join")

		c.Notify(*dto.NewMessageTransfer(dto.MessageFirstJoin, c.Game, response))
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
		err := c.Hub.GameService.StartGame(c.ConnectionCtx.Ctx, c.Game.Code)
		if err != nil {
			ResponseError(c, errors.Wrapf(err, "USER: %v | StartGame %s", c.User.Username, c.Game.Code))
			return
		}
		question, err := c.Hub.QuestionService.
			GetQuestion(c.ConnectionCtx.Ctx, c.Game.Code, uint(c.GameSocket.GameState.CurrentQuestionOffset))

		if err != nil {
			ResponseError(c, errors.Wrapf(err, "GetQuestion %s", c.Game.Code))
			return
		}
		err = c.Hub.QuestionService.Update(c.ConnectionCtx.Ctx, c.Game.Code, question)
		if err != nil {
			ResponseError(c, errors.Wrapf(err, "Update %s", c.Game.Code))
			return
		}

		c.GameSocket.GameState.SetQuestion(question)
		question.Start()
		c.GameSocket.GameState.Start()

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
		log.Info().Msgf("c.GameSocket.GameState: %v", c.GameSocket.GameState)
		if c.GameSocket.GameState.GameStage == models.GameStageShowQuestion {
			// For players
			for client := range c.GameSocket.ClientSet {
				userRank, err := c.Hub.LeaderboardService.GetUserRank(client.ConnectionCtx.Ctx, client.Game.Code, client.User.Username)

				if err != nil {
					ResponseError(c, errors.Wrapf(err, "USER %v | GetUserRank %s", client.User.Username, client.Game.Code))
					return
				}
				response := dto.MessageTransfer{
					Type: dto.MessageQuestionResult,
					Data: dto.NewQuestionResult(client.GameSocket.GameState.Question, client.QuestionAnswersMap[client.GameSocket.GameState.Question.ID], nil),
					Meta: map[string]interface{}{
						"userRank": userRank,
					},
				}
				client.Notify(response)
			}
			// For host
			questionStatistic, err := c.Hub.QuestionService.StatisticAnswersOfQuestion(c.ConnectionCtx.Ctx, c.Game.Code, uint(c.GameSocket.GameState.CurrentQuestionOffset))
			if err != nil {
				log.Error().Stack().Err(err).Msg("StatisticAnswersOfQuestion")
			}
			c.GameSocket.GameState.SetGameStage(models.GameStageShowAnswer)
			c.GameSocket.NotifyHost(dto.MessageTransfer{
				Type: dto.MessageQuestionResult,
				Data: dto.NewQuestionResult(c.GameSocket.GameState.Question, nil, questionStatistic),
			})
		} else {
			// else if c.GameSocket.GameState.GameStage == models.GameStageShowAnswer {
			c.GameSocket.GameState.SetGameStage(models.GameStageShowQuestion)
			c.GameSocket.GameState.NextQuestion()
			question, err := c.Hub.QuestionService.
				GetQuestion(c.ConnectionCtx.Ctx,
					c.Game.Code,
					uint(c.GameSocket.GameState.CurrentQuestionOffset))

			if err != nil {
				ResponseError(c, errors.Wrapf(err, "USER: %v | GetQuestion %s", c.User.Username, c.Game.Code))
				return
			}
			// It means game is ended
			if question == nil {
				c.GameSocket.GameState.SetStatus(models.GameStatusEnded)
				responseGameEnded(c)
				return
			}

			question.Start()
			err = c.Hub.QuestionService.Update(c.ConnectionCtx.Ctx, c.Game.Code, question)
			if err != nil {
				ResponseError(c, errors.Wrapf(err, "USER: %v | Update %s", c.User.Username, c.Game.Code))
				return
			}

			c.GameSocket.GameState.SetQuestion(question)
			c.GameSocket.GameState.SetGameStage(models.GameStageShowQuestion)
			response := dto.MessageTransfer{
				Type: dto.MessageNewQuestion,
				Data: dto.NewQuestionRes(question),
			}
			c.GameSocket.NotifyAll(response)
		}
	case dto.MessageLeaderboard:
		leaderBoard, err := c.Hub.LeaderboardService.GetLeaderboard(c.ConnectionCtx.Ctx, c.Game.Code)

		if err != nil {
			ResponseError(c, errors.Wrapf(err, "USER %v | GetLeaderboard %s", c.User.Username, c.Game.Code))
			return
		}

		c.Send <- dto.MessageTransfer{
			Type: dto.MessageLeaderboard,
			Data: leaderBoard.Items,
		}
	case dto.MessagePlayAgain:
		err := c.Hub.GameService.PlayAgain(c.ConnectionCtx.Ctx, c.Game.Code)

		if err != nil {
			ResponseError(c, errors.Wrapf(err, "PlayAgain %s", c.Game.Code))
			return
		}

		game, err := c.Hub.GameService.GetGame(c.ConnectionCtx.Ctx, c.Game.Code)
		if err != nil {
			ResponseError(c, errors.Wrapf(err, "GetGame %s", c.Game.Code))
			return
		}

		c.Game = game
		c.GameSocket.GameState.Reset()
		message := dto.MessageTransfer{
			Data: c.Game,
			Type: dto.MessageResetGame,
		}
		c.GameSocket.NotifyAll(message)
		c.GameSocket.NotifyUpdatedListPlayers()
	}
}

func playerMessageHandler(c *Client, message *dto.MessageTransfer) {
	log.Info().Msgf("playerMessageHandler - MessageTransfer: %v", message)
	switch message.Type {
	case dto.MessageAnswerQuestion:
		var answerDto dto.AnswerQuestionReq
		message.Binding(&answerDto)

		log.Info().Msgf("answerDto: %v", answerDto)
		answerDto.QuestionOffset = uint(c.GameSocket.GameState.CurrentQuestionOffset)
		answer, err := c.Hub.QuestionService.AnwserQuestion(c.ConnectionCtx.Ctx, c.Game.Code, c.User, answerDto)
		if err != nil {
			ResponseError(c, errors.Wrapf(err, "AnwserQuestion %v", answerDto))
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

func responseUserRankToPlayers(hub *Hub, clients map[*Client]bool) {
	for client := range clients {
		userRank, err := hub.LeaderboardService.GetUserRank(client.ConnectionCtx.Ctx, client.Game.Code, client.User.Username)
		if err != nil {
			ResponseError(client, errors.Wrapf(err, "USER %v | GetUserRank %s", client.User.Username, client.Game.Code))
			return
		}
		client.Notify(dto.MessageTransfer{
			Type: dto.MessageUserRank,
			Data: userRank,
		})
	}
}

func responseEndGameToPlayers(hub *Hub, clients map[*Client]bool) {
	for client := range clients {
		client.NotifyUserRank(dto.MessageEndGame)
	}
}

func responseGameEnded(host *Client) {
	host.GameSocket.GameState.SetStatus(models.GameStatusEnded)
	responseEndGameToPlayers(host.Hub, host.GameSocket.ClientSet)
	host.NotifyLeaderBoard(dto.MessageEndGame)
}
