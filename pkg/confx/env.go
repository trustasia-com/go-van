// Package confx provides ...
package confx

import (
	"os"
	"regexp"
)

var envValuePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(:-([^}]*))?\}`)

// ExpandEnv replaces ${VAR} and ${VAR:-default_value} expressions with environment values.
func ExpandEnv(data []byte) []byte {
	return []byte(expandEnvValue(string(data), os.LookupEnv))
}

func expandEnvValue(value string, lookupEnv func(string) (string, bool)) string {
	return envValuePattern.ReplaceAllStringFunc(value, func(expression string) string {
		parts := envValuePattern.FindStringSubmatch(expression)
		resolved, ok := lookupEnv(parts[1])
		if parts[2] != "" && (!ok || resolved == "") {
			return parts[3]
		}
		return resolved
	})
}
