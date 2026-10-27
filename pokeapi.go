package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rwndy/pokedex/internal/pokecache"
)

const locationAreasURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

const maxResponseBytes = 5 << 20 // 5 MiB

type locationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokeapiClient struct {
	httpClient *http.Client
	cache      *pokecache.Cache
}

func newPokeapiClient(timeout, cacheInterval time.Duration) *pokeapiClient {
	return &pokeapiClient{
		httpClient: &http.Client{Timeout: timeout},
		cache:      pokecache.NewCache(cacheInterval),
	}
}

func (c *pokeapiClient) Close() {
	c.cache.Close()
}

func (c *pokeapiClient) fetchLocationAreas(pageURL string) (locationAreasResponse, error) {
	body, err := c.get(pageURL)
	if err != nil {
		return locationAreasResponse{}, fmt.Errorf("fetching location areas: %w", err)
	}

	var page locationAreasResponse
	if err := json.Unmarshal(body, &page); err != nil {
		return locationAreasResponse{}, fmt.Errorf("decoding location areas: %w", err)
	}
	return page, nil
}

func (c *pokeapiClient) get(url string) ([]byte, error) {
	if body, ok := c.cache.Get(url); ok {
		return body, nil
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if len(body) > maxResponseBytes {
		return nil, fmt.Errorf("response larger than %d bytes", maxResponseBytes)
	}

	c.cache.Add(url, body)
	return body, nil
}
