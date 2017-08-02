package lanes

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/codekoala/go-aws-lanes/ssh"
)

type Profile struct {
	AWSAccessKeyId     string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`

	SSH ssh.Config `yaml:"ssh"`
}

// GetCurrentProfile loads the currently configured lane profile from the filesystem.
func (this *Config) GetCurrentProfile() *Profile {
	var (
		in   []byte
		prof = new(Profile)

		err error
	)

	ppath := this.GetProfilePath()

	if in, err = ioutil.ReadFile(ppath); err != nil {
		log.Fatalf("unable to read lane profile: %s", err)
	}

	if err = yaml.Unmarshal(in, prof); err != nil {
		log.Fatalf("unable to parse lane profile (%s): %s", ppath, err)
	}

	return prof
}
