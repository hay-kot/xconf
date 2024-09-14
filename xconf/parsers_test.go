package xconf

import "testing"

func TestParseArg(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		flags  []string
		expect string
	}{
		{
			name:   "Flag with equals sign",
			args:   []string{"--toml-file=path/to/file"},
			flags:  []string{"--toml-file"},
			expect: "path/to/file",
		},
		{
			name:   "Flag with space",
			args:   []string{"--toml-file", "path/to/file"},
			flags:  []string{"--toml-file"},
			expect: "path/to/file",
		},
		{
			name:   "Short flag with space",
			args:   []string{"-tf", "path/to/file"},
			flags:  []string{"-tf"},
			expect: "path/to/file",
		},
		{
			name:   "Multiple flags, first match",
			args:   []string{"--config", "config.yaml", "--toml-file", "path/to/file"},
			flags:  []string{"--toml-file", "--config"},
			expect: "config.yaml",
		},
		{
			name:   "Multiple flags, match second",
			args:   []string{"--config", "config.yaml", "--toml-file", "path/to/file"},
			flags:  []string{"--toml-file"},
			expect: "path/to/file",
		},
		{
			name:   "No match found",
			args:   []string{"--json-file", "data.json"},
			flags:  []string{"--toml-file"},
			expect: "",
		},
		{
			name:   "Flag with equals sign in middle of args",
			args:   []string{"--config=config.yaml", "--toml-file=path/to/file"},
			flags:  []string{"--toml-file"},
			expect: "path/to/file",
		},
		{
			name:   "Flag with no value",
			args:   []string{"--toml-file"},
			flags:  []string{"--toml-file"},
			expect: "",
		},
		{
			name:   "flag with short form",
			args:   []string{"-tf", "path/to/file"},
			flags:  []string{"--toml-file", "-tf"},
			expect: "path/to/file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArg(tt.args, tt.flags)
			if result != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, result)
			}
		})
	}
}
