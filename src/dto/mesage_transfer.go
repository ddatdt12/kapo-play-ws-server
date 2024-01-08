package dto

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type MessageType string

// for both
const (
	MessageFirstJoin      MessageType = "first_join"
	MessagePlayerList     MessageType = "player_list"
	MessageNewQuestion    MessageType = "new_question"
	MessageQuestionResult MessageType = "question_result"
	MessageEndGame        MessageType = "end_game"
	MessageResetGame      MessageType = "reset_game"
)

// For player
const (
	MessageAnswerQuestion MessageType = "player/answer_question"
	MessageUserRank       MessageType = "player/user_rank"
	MessagePlayerEndGame  MessageType = "player/end_game"
)

// For Host
const (
	MessageStartGame    MessageType = "host/start_game"
	MessageTimeUp       MessageType = "host/time_up"
	MessageSkipQuestion MessageType = "host/skip_game"
	MessageNextAction   MessageType = "host/next_action"
	MessageNextQuestion MessageType = "host/next_question"
	MessageLeaderboard  MessageType = "host/get_leaderboard"
	MessagePlayAgain    MessageType = "host/play_again"
)

// Error
const (
	Error MessageType = "error"
)

type MessageTransfer struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func NewMessageTransfer(messageType MessageType, data interface{}, meta interface{}) *MessageTransfer {
	return &MessageTransfer{
		Type: messageType,
		Data: data,
		Meta: meta,
	}
}

func (message *MessageTransfer) IsValid() bool {
	switch message.Type {
	case MessageFirstJoin,
		MessagePlayerList,
		MessageAnswerQuestion,
		MessageStartGame,
		MessageSkipQuestion,
		MessageNextQuestion,
		MessageLeaderboard,
		MessageTimeUp,
		MessageNextAction,
		MessageNewQuestion, MessageQuestionResult:
		return true
	default:
		log.Info().Interface("invalid messageType", message.Type).Msg("messageStruct")
		return false
	}
}

func (m *MessageTransfer) Binding(result any) error {
	return mapstructure.Decode(m.Data, result)
}
