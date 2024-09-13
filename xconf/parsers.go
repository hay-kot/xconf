package xconf

import (
	"strings"
)

// parseArg is a helper function to parse the arguments from the command line.
// It will return the first string match from the flags slice. If no match is
// found, it will return an empty string.
func parseArg(args []string, flags []string) string {
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Check for flags in --flag=value format
		for _, flag := range flags {
			if strings.HasPrefix(arg, flag+"=") {
				// Return the value after the "="
				return strings.SplitN(arg, "=", 2)[1]
			}

			// Check for --flag format, and the next argument as the value
			if arg == flag && i+1 < len(args) {
				return args[i+1]
			}
		}
	}

	// No match found
	return ""
}
