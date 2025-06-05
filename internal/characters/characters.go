package characters

import (
	"encoding/csv"
	"os"
)

type Character struct {
	Symbol   string
	Spelling string
	System   string
}

func Load(filename string) []Character {
	file, err := os.Open(filename)
	if err != nil {
		panic("Did not find library file")
	}
	defer file.Close()

	var characters []Character
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	for i, v := range records {
		if i != 0 { // TOOD: sure there's better way
			characters = append(characters, Character{Spelling: v[1], System: v[2], Symbol: v[0]})
		}
	}
	return characters
}
