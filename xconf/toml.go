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

// WithData accepts the toml document as a slice of bytes.
func WithData(data []byte) TOML {
	return TOML{
		data: data,
	}
}

// WithFileEnv accepts the environment variable names to look for the file path.
//
// The envs slice should contain the environment variable names to look for the
// file path. The first match will be used to extract the file path.
//
// If the file path is not found, an empty TOML struct is returned.
func WithFileEnv(envs []string) (TOML, error) {
	var filepath string
	for _, env := range envs {
		filepath = os.Getenv(env)
		if filepath != "" {
			break
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

// WithFileFlag accepts the command line arguments and flags to parse
// the file path from the command line arguments.
//
// The flags slice should contain the flag names to look for in the
// command line arguments. The first match will be used to extract the
// file path.
//
// If the flag is not found, an empty TOML struct is returned.
func WithFileFlag(args []string, flags ...string) (TOML, error) {
	filepath := parseArg(args, flags)
	if filepath == "" {
		return TOML{}, nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return TOML{}, fmt.Errorf("open file: %w", err)
	}

	return WithReader(f), nil
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
