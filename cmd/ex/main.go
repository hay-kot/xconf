package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/hay-kot/xconf/xconf"
)

type Nested struct {
	Name     string `toml:"name"      conf:"help:Name of the nested struct,notzero"`
	Number   int    `toml:"number"`
	UserPath string `toml:"user_path" xconf:"resolve"`
	RelPath  string `toml:"rel_path"  xconf:"resolve"`
}

type Config struct {
	LogLevel     string `toml:"log_level"     conf:"help:Log level"`
	EnableConfig bool   `toml:"enable_config" conf:"help:Enable the config"`
	// Nested is a nested struct
	Nested Nested `toml:"nested"`
}

const Prefix = "EX"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	cfg := &Config{}

	tomlProvider, err := xconf.WithFileSources(os.Args, []string{"--toml-file", "-tf"}, []string{"EX_TOML_FILE"})
	if err != nil {
		return err
	}

	help, err := conf.Parse(Prefix, cfg, tomlProvider)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	err = xconf.ResolvePaths(tomlProvider.FilePath(), cfg)
	if err != nil {
		return err
	}

	v, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(v))

	return nil
}
