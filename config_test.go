package lanes_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/codekoala/go-aws-lanes"
)

func TestLoadConfigBytes(t *testing.T) {
	envVar := "LANES_PROFILE"
	in := []byte(`profile: foo`)

	os.Unsetenv(envVar)
	out, err := lanes.LoadConfigBytes([]byte{})
	assert.NotNil(t, err)

	out, err = lanes.LoadConfigBytes(in)
	assert.Nil(t, err)
	assert.Equal(t, out.Profile, "foo")

	os.Setenv(envVar, "test")
	out, err = lanes.LoadConfigBytes(in)
	assert.Nil(t, err)
	assert.Equal(t, out.Profile, "test")
}
