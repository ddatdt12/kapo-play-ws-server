package dto

import (
	"github.com/rs/zerolog/log"
)

type MessageType string

const (
	// MessageTransferTypeMessage is the type of message transfer
	SendMessage MessageType = "send_message"
	// MessageTransferTypeUser is the type of message transfer
	NewMessage MessageType = "new_message"
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
