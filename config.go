package lanes

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

var (
	// CONFIG_DIR is the directory where all Lanes configuration files are expected to exist.
	CONFIG_DIR = EnvDefault("LANES_CONFIG_DIR", "$HOME/.lanes")

	// CONFIG is the path to the Lanes configuration file to use.
	CONFIG = EnvDefault("LANES_CONFIG", "$LANES_CONFIG_DIR/lanes.yml")

	// REGION is the default region to use for any profile that doesn't have a region explicitly set.
	REGION = EnvDefault("LANES_REGION", "us-west-2")

	// DISABLE_UTF8 switches the output tables to use ASCII instead of UTF-8 for borders.
	DISABLE_UTF8 = os.Getenv("LANES_DISABLE_UTF8") != ""

	// LANE_TAG is the name of the EC2 tag to use when determining which lane a server belongs in.
	LANE_TAG = EnvDefault("LANES_LANE_TAG", "Lane")

	// NAME_TAG is the name of the EC2 tag to use when determining a server's name.
	NAME_TAG = EnvDefault("LANES_NAME_TAG", "Name")
)

type Config struct {
	Profile string `yaml:"profile"`
}

// LoadConfig unmarshals the default YAML configuration file.
func LoadConfig() (*Config, error) {
	return LoadConfigFile(CONFIG)
}

// LoadConfigFile unmarshals the specified YAML file and returns a *Config.
func LoadConfigFile(cfgPath string) (c *Config, err error) {
	var in []byte

	if in, err = ioutil.ReadFile(cfgPath); err != nil {
		err = fmt.Errorf("unable to read configuration file: %s", err)
		return
	}

	return LoadConfigBytes(in)
}

// LoadConfigBytes unmarshals YAML input and returns a *Config with any environment variables taking precedence.
func LoadConfigBytes(in []byte) (c *Config, err error) {
	c = new(Config)

	if err = yaml.Unmarshal(in, c); err != nil {
		err = fmt.Errorf("unable to parse configuration: %s", err)
		return
	}

	// check the env for a profile, giving it precedence if set
	if envProfile := os.Getenv("LANES_PROFILE"); envProfile != "" {
		c.Profile = envProfile
	}

	// return an error if no profile is set in the environment or the specified config file
	if c.Profile == "" {
		err = fmt.Errorf("no profile specified; please specify it in %s or set LANES_PROFILE in your environment", CONFIG)
		return
	}

	return c, nil
}

// Write saves the current settings to the default configuration file.
func (this *Config) Write() (err error) {
	return this.WriteFile(CONFIG)
}

// WriteFile saves the current settings to the specified file.
func (this *Config) WriteFile(dest string) (err error) {
	var out []byte

	if out, err = this.WriteBytes(); err != nil {
		return
	}

	// make sure the destination directory exists
	if err = os.MkdirAll(path.Dir(dest), 0700); err != nil {
		return
	}

	return ioutil.WriteFile(dest, out, 0644)
}

// WriteBytes marshals the current settings to YAML.
func (this *Config) WriteBytes() ([]byte, error) {
	return yaml.Marshal(this)
}

// GetProfilePath determines where the current Lanes profile configuration file should be found.
func (this *Config) GetProfilePath() string {
	return path.Join(CONFIG_DIR, this.Profile+".yml")
}

// SetProfile changes the desired profile.
func (this *Config) SetProfile(name string) error {
	this.Profile = name
	return this.Write()
}
