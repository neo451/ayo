package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neo451/alpha/app"
	"github.com/neo451/alpha/internal/characters"
	"github.com/neo451/alpha/internal/config"
)

func loadLibrary(filename string) ([]characters.Character, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var characters []characters.Character
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&characters); err != nil {
		return nil, err
	}
	return characters, nil
}

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		fmt.Printf("Error loading config %v\n", err)
		return
	}

	characters, err := loadLibrary(cfg.Lib[0])

	if err != nil {
		panic("no library")
	}

	loop.Loop(cfg, characters)
}
