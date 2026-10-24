package main

func main() {
	firstPage := locationAreasURL
	cfg := &config{
		commands:         getCommands(),
		nextLocationsURL: &firstPage,
	}
	startRepl(cfg)
}
