package models

import (
	"github.com/rs/zerolog/log"
)

type Message struct {
	GameCode string `json:"game_code"`
	Message  string `json:"message"`
}

func NewMessage(data interface{}) *Message {
	// Assume x.Data contains a map with the necessary fields
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Info().Interface("data", data).Msg("x.Data is not a map[string]interface{}")
		return nil
	}

	// Convert the dataMap to a Message struct
	var messageStruct Message

	messageStruct.GameCode = dataMap["game_code"].(string) // Assuming GameID is a float in the map
	messageStruct.Message = dataMap["message"].(string)

	log.Info().Interface("messageStruct", messageStruct).Msg("messageStruct")
	return &messageStruct
}
