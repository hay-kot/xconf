package xconf

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	tag        = "xconf"
	tagResolve = "resolve"
)

// ResolvePaths walks the struct v and checks for the tag "conflib:resolve".
// If the tag is found, it resolves the path relative to the configpath.
// The configpath is the path to the configuration file. Absolute paths
// are ignored.
//
// Example Resolutions:
//   - './files.txt' -> '/path/to/configdir/files.txt'
//   - '~/files.txt'  -> '/home/user/files.txt'
func ResolvePaths(configpath string, v any) error {
	// ensure configpath is absolute
	if !filepath.IsAbs(configpath) {
		// convert to absolute path
		var err error
		configpath, err = filepath.Abs(configpath)
		if err != nil {
			panic(err)
		}
	}

	configdir := filepath.Dir(configpath)
	return resolvePathsRecursive(configdir, reflect.ValueOf(v))
}

// resolvePathsRecursive is a helper function that recursively resolves paths
func resolvePathsRecursive(configdir string, val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure we're dealing with a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", val.Kind())
	}

	// Iterate through the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := val.Type().Field(i)

		// Check if the field has the "conflib" tag with value "resolve"
		if tagValue, ok := structField.Tag.Lookup(tag); ok && tagValue == tagResolve {
			if field.Kind() == reflect.String && field.CanSet() {
				path := field.String()

				if path == "" {
					continue
				}

				// If the path is not absolute, resolve it relative to the baseDir
				if !filepath.IsAbs(path) {
					if strings.HasPrefix(path, "~") {
						homedir, err := os.UserHomeDir()
						if err != nil {
							return fmt.Errorf("get user home dir: %w", err)
						}

						path = strings.Replace(path, "~", homedir, 1)
					} else {
						path = filepath.Join(configdir, path)
					}

					field.SetString(path)
				}
			}
		}

		// Recursively resolve nested structs
		if field.Kind() == reflect.Struct {
			return resolvePathsRecursive(configdir, field)
		}

		// If the field is a pointer to a struct, resolve its paths recursively
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			return resolvePathsRecursive(configdir, field)
		}
	}

	return nil
}
