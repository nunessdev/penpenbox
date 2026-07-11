package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nunessdev/penpenbox/internal/models"
)

type steamGame struct {
	AppId           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

type response struct {
	GameCount int         `json:"game_count"`
	Games     []steamGame `json:"games"`
}

type apiResponse struct {
	Response response `json:"response"`
}

func GetGames(apiKey string, steamID string) ([]models.Game, error) {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&include_appinfo=1&format=json", apiKey, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Steam API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Steam API response: %w", err)
	}

	var raw apiResponse
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Steam API response: %w", err)
	}

	games := make([]models.Game, 0, len(raw.Response.Games))
	for _, g := range raw.Response.Games {
		games = append(games, models.Game{
			AppID:    g.AppId,
			Title:    g.Name,
			Playtime: g.PlaytimeForever,
			Platform: "steam",
		})
	}
	return games, nil
}
