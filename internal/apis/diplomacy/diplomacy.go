package diplomacy

import (
	"bytes"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"encoding/json"
	"errors"
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

func (a *API) GetGame(id string) (*models.Game, error) {
	url := fmt.Sprintf("%s/v1/games/%s", a.baseURL, id)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New(res.Status)
	}

	var g models.Game
	if err = json.NewDecoder(res.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (a *API) GetBoard(id string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/games/%s/board", a.baseURL, id)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New(res.Status)
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

	url := fmt.Sprintf("%s/v1/games", a.baseURL)
	res, err := http.Post(url, "application/json", bytes.NewReader(body))

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var body models.Error
		if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, errors.New(body.Message)
	}

	var g models.Game
	if err = json.NewDecoder(res.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}
