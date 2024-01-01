package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type IGameClient interface {
	Start(gameID uint) error
	End(gameID uint) error
	PlayAgain(gameID uint) error
}

type GameClient struct {
	ServiceURL string
}

func NewGameClient() *GameClient {
	serviceURL := configs.EnvConfigs.ANSWER_SERVICE_URL

	return &GameClient{
		ServiceURL: serviceURL,
	}
}

func (c *GameClient) Start(gameID uint) error {
	// Convert games to JSON
	jsonData, err := json.Marshal(GameRequest{
		Status: "playing",
		ID:     gameID,
	})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/games/%v/state", c.ServiceURL, gameID)
	// Send POST request to c.ServiceURL
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	}
	sb := string(body)
	log.Info().Msgf("Response body: %s", sb)

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
func (c *GameClient) End(gameID uint) error {
	// Convert games to JSON
	jsonData, err := json.Marshal(GameRequest{
		Status: "ended",
		ID:     gameID,
	})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/games/%v/state", c.ServiceURL, gameID)
	// Send POST request to c.ServiceURL
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	}
	sb := string(body)
	log.Info().Msgf("Response body: %s", sb)

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *GameClient) PlayAgain(gameID uint) error {
	// Convert games to JSON

	url := fmt.Sprintf("%s/games/%v/play-again", c.ServiceURL, gameID)
	// Send POST request to c.ServiceURL
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	}
	sb := string(body)
	log.Info().Msgf("Response body: %s", sb)

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return errors.New("Something went wrong when play again")
	}

	return nil
}
