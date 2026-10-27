package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Lists the Pokemon found in a location area",
			callback:    commandExplore,
		},
	}
}

func commandMap(cfg *config, args ...string) error {
	if cfg.nextLocationsURL == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return showLocationAreas(cfg, *cfg.nextLocationsURL)
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return showLocationAreas(cfg, *cfg.prevLocationsURL)
}

func showLocationAreas(cfg *config, pageURL string) error {
	page, err := cfg.pokeapiClient.fetchLocationAreas(pageURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = page.Next
	cfg.prevLocationsURL = page.Previous

	for _, area := range page.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("usage: explore <area_name>")
	}
	name := args[0]

	area, err := cfg.pokeapiClient.fetchLocationArea(name)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", name)
	fmt.Println("Found Pokemon:")
	for _, enc := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	names := make([]string, 0, len(cfg.commands))
	for name := range cfg.commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := cfg.commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
