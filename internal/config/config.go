package config

import (
	"github.com/BurntSushi/toml"
)

type PromptConfig struct {
	Format string `toml:"format"`
	Ok     string `toml:"ok"`
	Err    string `toml:"err"`
}

type CmdConfig struct {
	Exit string `toml:"exit"`
}

type ProgressConfig struct {
	Frequency int  `toml:"frequency"`
	Enabled   bool `toml:"enabled"`
}

type Config struct {
	Lib      []string       `toml:"lib"`
	Cmd      CmdConfig      `toml:"cmd"`
	Progress ProgressConfig `toml:"progress"`
	Prompt   PromptConfig   `toml:"prompt"`
	Theme    Theme
}

func DefaultConfig() Config {
	return Config{
		Lib: []string{"characters.json"},
		Cmd: CmdConfig{
			Exit: "q",
		},
		Prompt: PromptConfig{
			Ok:     "✅ Correct!",
			Err:    "❌ Incorrect. The answer is '%s'",
			Format: "[{{.System}}] {{.Symbol}}",
		},
		Progress: ProgressConfig{
			Frequency: 5,
		},
		Theme: DefaultTheme(),
	}
}

func Load(filename string) (Config, error) {
	cfg := DefaultConfig()
	if _, err := toml.Decode(filename, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c Config) Print() string {
	bytes, err := toml.Marshal(c)
	if err != nil {
	}
	return string(bytes)
}
