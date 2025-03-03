package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func main() {
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		input := cleanInput(text)
		command, ok := commands[input[0]]
		if !ok {
			fmt.Println("Unkown command")
			continue
		}
		err := command.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error exiting")
}

func commandHelp() error {
	listOfCommands := ""
	for key, value := range commands {
		listOfCommands += fmt.Sprintf("%s: %s\n", key, value.description)
	}

	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n%s", listOfCommands)
	return nil
}

// 2100 curl https://pokeapi.co/api/v2/location/
// 2101 curl https://pokeapi.co/api/v2/location/?offset=0&limit=20
// 2102 curl https://pokeapi.co/api/v2/location/?offset=20&limit=20

var offset = 0
var limit = 20

type MapResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func commandMap() error {
	offsetStr := fmt.Sprintf("offset=%v", offset)
	limitStr := fmt.Sprintf("limit=%v", limit)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?%v&%v", offsetStr, limitStr)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var maps MapResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&maps)
	if err != nil {
		return err
	}

	for _, result := range maps.Results {
		fmt.Println(result.Name)
	}

	offset += 20
	return nil
}

func commanMapb() error {
	offset -= 20
	if offset < 0 {
		offset = 0
	}
	offsetStr := fmt.Sprintf("offset=%v", offset)
	limitStr := fmt.Sprintf("limit=%v", limit)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location/?%v%v", offsetStr, limitStr)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var maps MapResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&maps)
	if err != nil {
		return err
	}

	for _, result := range maps.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
