package lanes_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/codekoala/go-aws-lanes"
)

func TestEnvDefault(t *testing.T) {
	name := "FOO"
	def := "default value"

	os.Unsetenv(name)
	assert.Equal(t, lanes.EnvDefault(name, def), def)
	assert.Equal(t, os.Getenv(name), def)

	os.Setenv(name, "custom")
	assert.Equal(t, lanes.EnvDefault(name, def), "custom")
}
