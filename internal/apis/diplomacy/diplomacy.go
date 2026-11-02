package diplomacy

import (
	"bytes"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	baseURL string
}

func NewAPI() API {
	return API{
		// baseURL: "https://orq2shzqa2.execute-api.us-west-2.amazonaws.com",
		baseURL: "http://localhost:8080",
	}
}

func (a *API) GetGame(gameID string) (*models.Game, error) {
	URL := fmt.Sprintf("%s/v1/games/%s", a.baseURL, gameID)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", res.Status)
	}

	var g models.Game
	if err = json.NewDecoder(res.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (a *API) GetBoard(gameID string) ([]byte, error) {
	URL := fmt.Sprintf("%s/v1/games/%s/board", a.baseURL, gameID)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) CreateGame(req models.CreateGameRequest) (*models.Game, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	URL := fmt.Sprintf("%s/v1/games", a.baseURL)
	res, err := http.Post(URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, decodeError(res)
	}

	var g models.Game
	if err = json.NewDecoder(res.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (a *API) CreatePlayer(gameID string, req models.CreatePlayerRequest) (*models.Player, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	URL := fmt.Sprintf("%s/v1/games/%s/players", a.baseURL, gameID)
	res, err := http.Post(URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, decodeError(res)
	}

	var p models.Player
	if err = json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (a *API) DeletePlayer(gameID string, playerID string) error {
	URL := fmt.Sprintf("%s/v1/games/%s/players/%s", a.baseURL, gameID, playerID)
	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return decodeError(res)
	}

	return nil
}

func (a *API) CreateOrder(gameID string, turnID string, req models.CreateOrderRequest) (models.CommandValidationResults, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	URL := fmt.Sprintf("%s/v1/games/%s/turns/%s/orders", a.baseURL, gameID, turnID)
	res, err := http.Post(URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, decodeError(res)
	}

	var results models.CommandValidationResults
	if err = json.NewDecoder(res.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}

func decodeError(res *http.Response) error {
	var body models.Error
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}
	if body.Message == "" {
		return fmt.Errorf("%s", res.Status)
	}

	return fmt.Errorf("%s", body.Message)
}
