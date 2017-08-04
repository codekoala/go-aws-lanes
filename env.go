package lanes

import (
	"os"
)

// EnvDefault returns the value of the specified environment variable. If that variable is not set, the first non-empty
// value will be returned instead. The returned value will also be set in the environment for later use.
func EnvDefault(varName string, values ...string) (value string) {
	values = append([]string{os.Getenv(varName)}, values...)
	for _, value = range values {
		if value != "" {
			value = os.ExpandEnv(value)

			// set the value in the environment for future use
			os.Setenv(varName, value)
		}
	}

	return value
}
