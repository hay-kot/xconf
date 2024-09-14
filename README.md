# xconf

xconf provides some extensions to the conf package for generating docs and a toml provider

[Go Reference](https://pkg.go.dev/github.com/hay-kot/xconf)

## Install

```bash
go get -u github.com/hay-kot/xconf
```

## Features

- TOML Provider for sourcing configuration from a TOML File.
- Read TOML path from cli args or environment variables.
- Path Resolver for resolving relative paths in configuration files
  - Resolves `./path` to be relative to the configuration file parent directory.
  - Resolves `~/path` to be relative to the user's home directory.
