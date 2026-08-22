package diplomacy

import (
	"discord-diplomacy/internal/apis/diplomacy/models"
	"encoding/json"
	"errors"
	"fmt"
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
