package answer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/rs/zerolog/log"
)

type IAnswerClient interface {
	CreateAnswer(answer *CreateAnswer) error
}

type AnswerClient struct {
	ServiceURL string
}

func NewAnswerClient() *AnswerClient {
	serviceURL := configs.EnvConfigs.ANSWER_SERVICE_URL

	return &AnswerClient{
		ServiceURL: serviceURL,
	}
}

const (
	TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYxOTg0MTQsIkRhdGEiOnsiZW1haWwiOiJ0ZXN0QGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxIn19.nEqv63rFygknkGbcpwuoWpQCI5XTnFigqhg8dxL8YYY"
)

func (c *AnswerClient) GetTemplates() error {
	url := fmt.Sprintf("%s/templates", c.ServiceURL)

	log.Info().Msgf("URL: %s", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TOKEN))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("Error sending request to server")
	}

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

func (c *AnswerClient) CreateAnswer(answer *CreateAnswer) error {
	// Convert answer to JSON
	jsonData, err := json.Marshal(answer)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/games/%v/answers", c.ServiceURL, answer.GameID)
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
