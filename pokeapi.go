package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const locationAreasURL = "https://pokeapi.co/api/v2/location-area/"

type locationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func fetchLocationAreas(pageURL string) (locationAreasResponse, error) {
	resp, err := httpClient.Get(pageURL)
	if err != nil {
		return locationAreasResponse{}, fmt.Errorf("fetching location areas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return locationAreasResponse{}, fmt.Errorf("fetching location areas: unexpected status %s", resp.Status)
	}

	var page locationAreasResponse
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return locationAreasResponse{}, fmt.Errorf("decoding location areas: %w", err)
	}
	return page, nil
}
