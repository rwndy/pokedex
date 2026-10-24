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
	}
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
