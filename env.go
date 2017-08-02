package lanes

import (
	"os"
)

// EnvDefault returns the value of the specified environment variable. If that variable is not set, the specified
// default value will be returned instead. The returned value will also be set in the environment for later use.
func EnvDefault(varName, defaultValue string) (value string) {
	value = os.Getenv(varName)
	if value == "" {
		value = os.ExpandEnv(defaultValue)

		// set the value in the environment for future use
		os.Setenv(varName, value)
	}

	return value
}
