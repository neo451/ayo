package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

func getXDGPath(envVar, defaultSubdir string) string {
	base := os.Getenv(envVar)
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("unable to determine home directory")
		}
		base = filepath.Join(home, defaultSubdir)
	}
	return filepath.Join(base, "alpha")
}

// TODO: don't copy if data is there
func moveProjectDataToXDGData() error {
	// Resolve paths
	projectDataDir := "data" // relative to where app is run from
	dataDir := getXDGPath("XDG_DATA_HOME", ".local/share")

	// Ensure target data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// Walk through files in ./data
	return filepath.Walk(projectDataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		destPath := filepath.Join(dataDir, info.Name())

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		fmt.Printf("Copied %s → %s\n", path, destPath)
		return nil
	})
}

func main() {

	configDir := getXDGPath("XDG_CONFIG_HOME", ".config")
	dataDir := getXDGPath("XDG_DATA_HOME", ".local/share")

	os.MkdirAll(configDir, 0755)
	os.MkdirAll(dataDir, 0755)

	configFile := filepath.Join(configDir, "config.toml")

	_, config_err := os.ReadFile(configFile)
	if config_err != nil {
		os.WriteFile(configFile, []byte("# **some helpful comments**\n"), 0644)
	}

	if err := moveProjectDataToXDGData(); err != nil {
		fmt.Println("Error moving data:", err)
	}

	cfg, err := config.Load(filepath.Join(configDir, "config.toml"))
	if err != nil {
		fmt.Printf("Error loading config %v\n", err)
		return
	}

	characters, err := loadLibrary(filepath.Join(dataDir, cfg.Lib[0]))

	if err != nil {
		panic("no library loaded")
	}

	// _ = loop.Loop
	// _ = characters
	loop.Loop(cfg, characters)
}
