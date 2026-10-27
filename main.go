package main

import "time"

func main() {
	client := newPokeapiClient(10*time.Second, 5*time.Minute)
	defer client.Close()

	firstPage := locationAreasURL
	cfg := &config{
		commands:         getCommands(),
		pokeapiClient:    client,
		nextLocationsURL: &firstPage,
		pokedex:          make(map[string]pokemon),
	}
	startRepl(cfg)
}
