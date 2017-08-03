package lanes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/codekoala/go-aws-lanes/ssh"
)

var (
	ErrMissingAccessKey = errors.New("missing AWS access key ID")
	ErrMissingSecretKey = errors.New("missing AWS secret key")
)

type Profile struct {
	AWSAccessKeyId     string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
	Region             string `yaml:"region"`

	SSH ssh.Config `yaml:"ssh"`
}

func (this *Profile) Validate() error {
	if this.AWSAccessKeyId == "" {
		return ErrMissingAccessKey
	}

	if this.AWSSecretAccessKey == "" {
		return ErrMissingSecretKey
	}

	if this.Region == "" {
		this.Region = REGION
	} else {
		REGION = this.Region
	}

	return nil
}

func (this *Profile) Activate() {
	os.Setenv("AWS_ACCESS_KEY_ID", this.AWSAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", this.AWSSecretAccessKey)
}

func (this *Profile) Deactivate() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
}

// GetCurrentProfile loads the currently configured lane profile from the filesystem.
func (this *Config) GetCurrentProfile() (prof *Profile, err error) {
	var in []byte

	ppath := this.GetProfilePath()

	if in, err = ioutil.ReadFile(ppath); err != nil {
		err = fmt.Errorf("unable to read lane profile: %s", err)
		return
	}

	prof = new(Profile)
	if err = yaml.Unmarshal(in, prof); err != nil {
		err = fmt.Errorf("unable to parse lane profile (%s): %s", ppath, err)
		return
	}

	if err = prof.Validate(); err != nil {
		err = fmt.Errorf("invalid profile: %s", err)
		return
	}

	return prof, nil
}
