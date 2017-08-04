package lanes

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/yaml.v2"

	"github.com/codekoala/go-aws-lanes/ssh"
)

type Profile struct {
	AWSAccessKeyId     string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
	Region             string `yaml:"region"`

	SSH ssh.Config `yaml:"ssh"`

	global *Config
}

// Validate checks that the profile includes the necessary information to interact with AWS.
func (this *Profile) Validate() error {
	if this.AWSAccessKeyId == "" {
		return ErrMissingAccessKey
	}

	if this.AWSSecretAccessKey == "" {
		return ErrMissingSecretKey
	}

	if this.global != nil {
		if this.Region == "" {
			this.Region = this.global.Region
		}
	} else {
		this.Region = os.Getenv("LANES_REGION")
	}

	return nil
}

// Activate sets some environment variables to access AWS using a given profile.
func (this *Profile) Activate() {
	os.Setenv("AWS_ACCESS_KEY_ID", this.AWSAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", this.AWSSecretAccessKey)
}

// Deactivate unsets environment variables to no longer access AWS with this profile.
func (this *Profile) Deactivate() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
}

// FetchServers retrieves all EC2 instances for the current profile.
func (this *Profile) FetchServers(svc *ec2.EC2) ([]*Server, error) {
	return this.FetchServersBy(svc, nil)
}

// FetchServersInLane retrieves all EC2 instances in the specified lane for the current profile.
func (this *Profile) FetchServersInLane(svc *ec2.EC2, lane string) ([]*Server, error) {
	return this.FetchServersBy(svc, CreateLaneFilter(lane))
}

// FetchServersBy retrieves all EC2 instances for the current profile using any specified filters. Each instance is
// automatically tagged with the appropriate SSH profile to access it.
func (this *Profile) FetchServersBy(svc *ec2.EC2, input *ec2.DescribeInstancesInput) (servers []*Server, err error) {
	var exists bool

	if servers, err = FetchServersBy(svc, input); err != nil {
		return
	}

	for _, svr := range servers {
		if svr.profile, exists = this.SSH.Mods[svr.Lane]; !exists {
			fmt.Printf("WARNING: no profile found for server %s", svr)
		}
	}

	return servers, nil
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

	// allow the profile to access global configuration values
	prof.global = this

	if err = prof.Validate(); err != nil {
		err = fmt.Errorf("invalid profile: %s", err)
		return
	}

	return prof, nil
}
