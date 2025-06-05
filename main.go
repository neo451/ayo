package main

import (
	"fmt"
	"github.com/neo451/ayo/app"
	char "github.com/neo451/ayo/internal/characters"
	"github.com/neo451/ayo/internal/config"
	"io"
	"os"
	"path/filepath"
	// "github.com/neo451/ayo/app/stat"
)

func getXDGPath(envVar, defaultSubdir string) string {
	base := os.Getenv(envVar)
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("unable to determine home directory")
		}
		base = filepath.Join(home, defaultSubdir)
	}
	return filepath.Join(base, "ayo")
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
		return nil
	})
}

func setupFiles() (string, string) {
	configDir := getXDGPath("XDG_CONFIG_HOME", ".config")
	dataDir := getXDGPath("XDG_DATA_HOME", ".local/share")

	os.MkdirAll(configDir, 0755)
	os.MkdirAll(dataDir, 0755)

	configFile := filepath.Join(configDir, "config.toml")

	_, configErr := os.ReadFile(configFile)
	if configErr != nil {
		os.WriteFile(configFile, []byte("# **some helpful comments**\n"), 0644)
	}

	if err := moveProjectDataToXDGData(); err != nil {
		fmt.Println("Error moving data:", err)
	}
	return configDir, dataDir
}

// Init data and config
func setup() config.Config {
	configDir, dataDir := setupFiles()
	configPath := filepath.Join(configDir, "config.toml")
	configStr, read_err := os.ReadFile(configPath)

	if read_err != nil {
		panic("no config file")
	}

	cfg, lib_err := config.Load(string(configStr))
	if lib_err != nil {
		panic(fmt.Sprintf("Error loading config %v\n", lib_err))
	}
	cfg.DataDir = dataDir
	return cfg
}

func main() {
	cfg := setup()
	characters := char.Load(filepath.Join(cfg.DataDir, cfg.Lib[0]))

	app.Quiz(cfg, characters)
	// app.Card(cfg, characters)
	// stat.RenderStat(characters)
}
