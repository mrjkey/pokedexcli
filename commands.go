package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

var commands map[string]cliCommand

func initCommands() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    commanMapb,
		},
		"explore": {
			name:        "explore",
			description: "Show encounters in a specified area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Pokedex listed",
			callback:    commandPokedex,
		},
	}
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error exiting")
}

func commandHelp(args []string) error {
	listOfCommands := ""
	for key, value := range commands {
		listOfCommands += fmt.Sprintf("%s: %s\n", key, value.description)
	}

	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n%s", listOfCommands)
	return nil
}

func commandExplore(args []string) error {
	if len(args) < 2 {
		return errors.New("missing location area name")
	}
	// fmt.Println(args[1])
	locationAreaName := args[1]

	content, err := getRequestWithName(locationAreaName)
	if err != nil {
		return err
	}

	err = printPokemon(content)
	if err != nil {
		return err
	}

	return nil
}

func commandMap(args []string) error {
	body, err := getRequest()
	if err != nil {
		return err
	}

	err = printMaps(body)
	if err != nil {
		return err
	}

	offset += 20
	return nil
}

func commanMapb(args []string) error {
	offset -= 20
	if offset < 0 {
		offset = 0
	}

	body, err := getRequest()
	if err != nil {
		return err
	}

	err = printMaps(body)
	if err != nil {
		return err
	}

	return nil
}

func commandCatch(args []string) error {
	if len(args) < 2 {
		return errors.New("missing pokemon name")
	}

	pokemonName := args[1]

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	content, err := getPokemonRequest(pokemonName)
	if err != nil {
		return err
	}

	var pokemon Pokemon
	err = json.Unmarshal(content, &pokemon)
	if err != nil {
		return err
	}

	baseXP := pokemon.BaseExperience
	random := rand.Intn(baseXP)
	// fmt.Printf("%v's base xp: %v\n", pokemonName, baseXP)

	if random < 20 {
		fmt.Printf("%v was caught!\n", pokemonName)
		pokedex[pokemonName] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(args []string) error {
	if len(args) < 2 {
		return errors.New("no name given")
	}

	pokemon, ok := pokedex[args[1]]
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	printPokemonStats(pokemon)

	return nil
}

func commandPokedex(args []string) error {
	fmt.Println("Your Pokedex:")
	for key := range pokedex {
		fmt.Printf(" - %v\n", key)
	}
	return nil
}
