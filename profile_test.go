package lanes_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/codekoala/go-aws-lanes"
)

func TestProfileValidate(t *testing.T) {
	p := &lanes.Profile{}
	assert.Equal(t, p.Region, "")

	assert.NotNil(t, p.Validate())

	p.AWSAccessKeyId = "demo"
	assert.NotNil(t, p.Validate())

	p.AWSSecretAccessKey = "demo"
	assert.Nil(t, p.Validate())

	// default region set by env var
	assert.Equal(t, p.Region, lanes.REGION)
}

func TestProfileActivate(t *testing.T) {
	key := "AWS_ACCESS_KEY_ID"
	secret := "AWS_SECRET_ACCESS_KEY"
	p := &lanes.Profile{
		AWSAccessKeyId:     "foo",
		AWSSecretAccessKey: "bar",
	}

	p.Activate()
	assert.Equal(t, os.Getenv(key), "foo")
	assert.Equal(t, os.Getenv(secret), "bar")

	p.Deactivate()
	assert.Equal(t, os.Getenv(key), "")
	assert.Equal(t, os.Getenv(secret), "")
}
