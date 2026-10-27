package main

import (
	"fmt"
	"os"
	"sort"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandMap(cfg *config) error {
	if cfg.nextLocationsURL == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return showLocationAreas(cfg, *cfg.nextLocationsURL)
}

func commandMapb(cfg *config) error {
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
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
