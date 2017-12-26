package lanes

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"

	"github.com/codekoala/go-aws-lanes/ssh"
)

type Profile struct {
	AWSAccessKeyId     string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
	Region             string `yaml:"region,omitempty"`

	SSH ssh.Config `yaml:"ssh"`

	global    *Config
	overwrite bool
}

// GetSampleProfile returns a sample profile that is easy to use as an example.
func GetSampleProfile() *Profile {
	return &Profile{
		SSH: ssh.Config{
			Default: &ssh.DefaultProfile,
			Mods: map[string]*ssh.Profile{
				"dev": {
					Identity: "~/.ssh/id_rsa_dev",
					Tunnels: []string{
						"8080:127.0.0.1:80",
						"3306:127.0.0.1:3306",
					},
				},
				"stage": {
					Identity: "~/.ssh/id_rsa_stage",
					Tunnel:   "8080:127.0.0.1:80",
				},
				"prod": {
					Identity: "~/.ssh/id_rsa_prod",
				},
			},
		},
	}
}

// GetAvailableProfiles returns a list of all Lanes profiles found in the configuration directory.
func GetAvailableProfiles() (found []string) {
	matches, _ := filepath.Glob(filepath.Join(ConfigDir, "*.yml"))
	for _, match := range matches {
		name := filepath.Base(strings.TrimSuffix(match, filepath.Ext(match)))
		if name == "lanes" {
			// this is not a profile
			continue
		}

		found = append(found, name)
	}

	return
}

// GetProfilePath uses the specified name to return a path to the file that is expected to hold the configuration for
// the named profile.
func GetProfilePath(name string, checkPerms bool) string {
	path := filepath.Join(ConfigDir, name+".yml")

	if checkPerms {
		CheckProfilePermissions(path)
	}

	return path
}

// CheckProfilePermissions looks for any concerns with permissions that are too permissible for Lanes profiles.
func CheckProfilePermissions(path string) {
	var result error

	// check the directory first
	dFatal, dErr := CheckPermissions(filepath.Dir(path))
	if dErr != nil {
		result = multierror.Append(dErr)
	}

	// check the actual profile
	pFatal, pErr := CheckPermissions(path)
	if pErr != nil {
		result = multierror.Append(pErr)
	}

	prefix := "WARNING"
	fatal := dFatal || pFatal
	if fatal {
		prefix = "ERROR"
	}

	if result != nil {
		fmt.Printf("%s: checking profile permissions, %s\n\n", prefix, result)
	}

	if fatal {
		os.Exit(1)
	}
}

// CheckPermissions looks for possible concerns with directory and file permissions.
func CheckPermissions(path string) (fatal bool, result error) {
	pStats, err := os.Stat(path)
	if err != nil {
		fatal = true
		result = multierror.Append(result, err)
	} else {
		mode := pStats.Mode()

		// check user permissions
		if (mode&0700)>>6 <= 3 {
			fatal = true
			result = multierror.Append(result, fmt.Errorf("%s is not user-accessible", path))
		}

		// check group permissions
		if (mode&0070)>>3 != 0 {
			result = multierror.Append(result, fmt.Errorf("%s is group-accessible", path))
		}

		// check world permissions
		if mode&0007 != 0 {
			result = multierror.Append(result, fmt.Errorf("%s is world-accessible", path))
		}
	}

	return
}

// LoadProfile attempts to read the specified profile from the filesystem.
func LoadProfile(name string) (prof *Profile, err error) {
	var in []byte

	if in, err = ioutil.ReadFile(GetProfilePath(name, true)); err != nil {
		err = fmt.Errorf("unable to read profile: %s", err)
		return
	}

	return LoadProfileBytes(in)
}

// LoadProfileBytes loads the currently configured lane profile from the specified YAML bytes.
func LoadProfileBytes(in []byte) (prof *Profile, err error) {
	prof = new(Profile)
	if err = yaml.Unmarshal(in, prof); err != nil {
		err = fmt.Errorf("unable to parse lane profile: %s", err)
		return
	}

	// allow the profile to access global configuration values
	prof.global = config

	if err = prof.Validate(); err != nil {
		err = fmt.Errorf("invalid profile: %s", err)
		return
	}

	return prof, nil
}

// SetOverwrite allows other packages to mark this profile as one that can safely be overwritten.
func (this *Profile) SetOverwrite(value bool) {
	this.overwrite = value
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
	return this.FetchServersBy(svc, nil, "")
}

// FetchServersInLane retrieves all EC2 instances in the specified lane for the current profile.
func (this *Profile) FetchServersInLane(svc *ec2.EC2, lane string) ([]*Server, error) {
	return this.FetchServersInLaneByKeyword(svc, lane, "")
}

// FetchServersInLane retrieves all EC2 instances in the specified lane for the current profile.
func (this *Profile) FetchServersInLaneByKeyword(svc *ec2.EC2, lane, keyword string) ([]*Server, error) {
	return this.FetchServersBy(svc, CreateLaneFilter(lane), keyword)
}

// FetchServersBy retrieves all EC2 instances for the current profile using any specified filters. Each instance is
// automatically tagged with the appropriate SSH profile to access it.
func (this *Profile) FetchServersBy(svc *ec2.EC2, input *ec2.DescribeInstancesInput, keyword string) (servers []*Server, err error) {
	if servers, err = FetchServersBy(svc, input, keyword); err != nil {
		return
	}

	return servers, nil
}

// Write saves the current settings to disk using the specified profile name.
func (this *Profile) Write(name string) (err error) {
	return this.WriteFile(name, GetProfilePath(name, false))
}

// WriteFile saves the current profile settings to the specified file.
func (this *Profile) WriteFile(name, dest string) (err error) {
	var out []byte

	// don't overwrite existing profiles without a flag being set to allow it
	if _, err = os.Stat(dest); err == nil && !this.overwrite {
		return fmt.Errorf("profile %q already exists", name)
	}

	if out, err = this.WriteBytes(); err != nil {
		return
	}

	// make sure the destination directory exists
	if err = os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
		return
	}

	if err = ioutil.WriteFile(dest, out, 0600); err != nil {
		return
	}

	fmt.Printf("Profile %q written to %s\n", name, dest)

	return nil
}

// WriteBytes marshals the current settings to YAML.
func (this *Profile) WriteBytes() ([]byte, error) {
	return yaml.Marshal(this)
}
