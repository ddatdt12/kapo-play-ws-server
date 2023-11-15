package dto

import (
	"github.com/rs/zerolog/log"
)

type MessageType string

const (
	// MessageTransferTypeSendMessage is the type of message transfer
	SendMessage    MessageType = "send_message"
	NewMessage     MessageType = "new_message"
	StartGame      MessageType = "start_game"
	SkipGame       MessageType = "skip_game"
	NextQuestion   MessageType = "next_question"
	AnswerQuestion MessageType = "answer_question"
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

func VerifyMessageType(messageType MessageType) bool {
	switch messageType {
	case SendMessage, NewMessage:
		return true
	default:
		log.Info().Interface("invalid messageType", messageType).Msg("messageStruct")
		return false
	}
}
