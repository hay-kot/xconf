package xconf

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// TOML provides support for unmarshaling TOML into the applications
// config value. After the toml is unmarshaled, the Parse function is
// executed to apply defaults and overrides. Fields that are not set to
// their zero after the toml is parsed will have the defaults ignored.
type TOML struct {
	data []byte
}

// WithFileSources accepts a list of flags and envs to search for the
// file path to read the toml from. The first file found will be used.
// If no file is found, an empty TOML is returned.
//
// Example:
//
//	cfg, err := xconf.WithFileSources([]string{"--config", "-c"}, []string{"CONFIG_FILE"})
//	if err != nil {
//		log.Fatal(err)
//	}
func WithFileSources(flags []string, envs []string) (TOML, error) {
	var filepath string
	filepath = parseArg(os.Args, flags)
	if filepath == "" {
		for _, env := range envs {
			filepath = os.Getenv(env)
			if filepath != "" {
				break
			}
		}
	}

	if filepath == "" {
		return TOML{}, nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return TOML{}, fmt.Errorf("open file: %w", err)
	}

	return WithReader(f), nil
}

func WithData(data []byte) TOML {
	return TOML{data: data}
}

// WithReader accepts a reader to read the toml.
func WithReader(r io.Reader) TOML {
	var b bytes.Buffer
	if _, err := b.ReadFrom(r); err != nil {
		return TOML{}
	}

	return WithData(b.Bytes())
}

// Process performs the actual processing of the toml.
func (y TOML) Process(prefix string, cfg interface{}) error {
	err := toml.Unmarshal(y.data, cfg)
	if err != nil {
		return fmt.Errorf("unmarshal toml: %w", err)
	}
	return nil
}
