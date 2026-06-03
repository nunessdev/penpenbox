package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Game struct {
	AppId           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

type Response struct {
	GameCount int    `json:"game_count"`
	Games     []Game `json:"games"`
}

type APIResponse struct {
	Response Response `json:"response"`
}

func GetGames(apiKey string, steamID string) ([]Game, error) {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&include_appinfo=1&format=json", apiKey, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Steam API: %w\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Steam API response: %w\n", err)
	}

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Steam API response: %w\n", err)
	}

	return apiResponse.Response.Games, nil
}
