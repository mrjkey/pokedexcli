package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mrjkey/pokedexcli/internal/pokecache"
)

var cache pokecache.Cache

func main() {
	initCommands()
	// rand.Seed(time.Now().UnixNano())
	cache = pokecache.NewCache(time.Second * 5)
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
		err := command.callback(input)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 2100 curl https://pokeapi.co/api/v2/location/
// 2101 curl https://pokeapi.co/api/v2/location/?offset=0&limit=20
// 2102 curl https://pokeapi.co/api/v2/location/?offset=20&limit=20

var offset = 0
var limit = 20

func getRequest() ([]byte, error) {
	offsetStr := fmt.Sprintf("offset=%v", offset)
	limitStr := fmt.Sprintf("limit=%v", limit)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?%v&%v", offsetStr, limitStr)

	return processRequest(url)
}

func getPokemonRequest(pokemonName string) ([]byte, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", pokemonName)
	return processRequest(url)
}

func processRequest(url string) ([]byte, error) {
	body, ok := cache.Get(url)
	if ok {
		return body, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	cache.Add(url, body)

	return body, nil
}

func getRequestWithName(name string) ([]byte, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", name)
	return processRequest(url)
}

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

type EncounterResponse struct {
	Encounters []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Url            string `json:"url"`
}

func printMaps(body []byte) error {
	var maps MapResponse
	err := json.Unmarshal(body, &maps)
	if err != nil {
		return err
	}

	for _, result := range maps.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func printPokemon(content []byte) error {
	var encResp EncounterResponse
	err := json.Unmarshal(content, &encResp)
	if err != nil {
		return err
	}

	for _, result := range encResp.Encounters {
		fmt.Println(result.Pokemon.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
