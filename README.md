# go-aws-lanes
[![Travis CI Status](https://travis-ci.org/codekoala/go-aws-lanes.svg?branch=master)](https://travis-ci.org/codekoala/go-aws-lanes)
[![License BSD3](https://img.shields.io/badge/license-BSD3-blue.svg)](https://raw.githubusercontent.com/codekoala/go-aws-lanes/master/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/codekoala/go-aws-lanes/total.svg)](https://github.com/codekoala/go-aws-lanes/releases)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/codekoala/go-aws-lanes)

This command line tool is designed to help people interact with different sets
of AWS profiles and EC2 instances. It allows you to easily switch between
multiple sets of AWS credentials and perform the following operations:

* list EC2 instances on the account, optionally filtered by a "Lane" tag.
* quickly SSH into a specific EC2 instance using the correct credentials,
  optionally setting up tunnels to locally access services running on a given
  instance.
* copy files to all EC2 instances in a given lane
* run commands on all EC2 instances in a given lane

## Sample Output

    $ lanes ls
    Current profile: foo
    Fetching servers... done
    ╭────────────────────────────────────────────────╮
    │                  AWS Servers                   │
    ├─────┬──────┬─────────┬────────────┬────────────┤
    │ IDX │ LANE │ SERVER  │ IP ADDRESS │ ID         │
    ├─────┼──────┼─────────┼────────────┼────────────┤
    │ 1   │ dev  │ dev-01  │ 1.2.3.4    │ i-12341234 │
    │ 2   │ uat  │ uat-01  │ 1.2.3.5    │ i-12341235 │
    │ 3   │ prod │ prod-01 │ 1.2.3.6    │ i-12341236 │
    │ 4   │ prod │ prod-02 │ 1.2.3.7    │ i-12341237 │
    ╰─────┴──────┴─────────┴────────────┴────────────╯

    $ lanes ls dev
    Current profile: foo
    Fetching servers... done
    ╭────────────────────────────────────────────────╮
    │                  AWS Servers                   │
    ├─────┬──────┬─────────┬────────────┬────────────┤
    │ IDX │ LANE │ SERVER  │ IP ADDRESS │ ID         │
    ├─────┼──────┼─────────┼────────────┼────────────┤
    │ 1   │ dev  │ dev-01  │ 1.2.3.4    │ i-12341234 │
    ╰─────┴──────┴─────────┴────────────┴────────────╯

    $ lanes ls prod
    Current profile: foo
    Fetching servers... done
    ╭────────────────────────────────────────────────╮
    │                  AWS Servers                   │
    ├─────┬──────┬─────────┬────────────┬────────────┤
    │ IDX │ LANE │ SERVER  │ IP ADDRESS │ ID         │
    ├─────┼──────┼─────────┼────────────┼────────────┤
    │ 1   │ prod │ prod-01 │ 1.2.3.6    │ i-12341236 │
    │ 2   │ prod │ prod-02 │ 1.2.3.7    │ i-12341237 │
    ╰─────┴──────┴─────────┴────────────┴────────────╯

## What Are Lanes?

A lane is basically a logical environment for your EC2 instances. For example,
you could have a lane called "dev" for development servers, one called "uat"
user acceptance testing, and one called "prod" for production servers.

## Installation

1. Download a pre-compiled, released version from [the releases
   page](https://github.com/codekoala/go-aws-lanes/releases)
2. Mark the binary as executable
3. Move the binary to your ``$PATH``

For example:

```shell
$ curl -Lo /tmp/lanes https://github.com/codekoala/go-aws-lanes/releases/download/v0.1.1/lanes_linux_amd64
$ chmod +x /tmp/lanes
$ sudo mv /tmp/lanes /usr/local/bin/lanes
```

To compile from source, please see the instructions in the [contributing section](#contributing).

## Configuration

The configuration for this tool lives in ``$HOME/.lanes/`` by default. Create a
``$HOME/.lanes/lanes.yml`` file with the following content:

```yaml
profile: demo
```

Next, create a ``$HOME/.lanes/demo.yml`` file with contents such as the
following:

```yaml
aws_access_key_id: [your AWS_ACCESS_KEY_ID for the "demo" profile]
aws_secret_access_key: [your AWS_SECRET_ACCESS_KEY for the "demo" profile]
ssh:
  mods:
    dev:
      identity: ~/.ssh/id_rsa_demo_dev
      tunnels:
        - 8080:127.0.0.1:80
        - 3306:127.0.0.1:3306
    uat:
      identity: ~/.ssh/id_rsa_demo_uat
      tunnel: 8080:127.0.0.1:80
    prod:
      identity: ~/.ssh/id_rsa_demo_prod
```

You can create additional profiles by creating new YAML files using this
pattern: ``$HOME/.lanes/[profile name].yml``

### Environment Variables

``lanes`` supports a handful of environment variables to quickly change
behavior:

* ``LANES_CONFIG_DIR``: the directory where all configuration is expected to
  reside. Default: ``$HOME/.lanes/``
* ``LANES_CONFIG``: the configuration file to use for lanes. Default:
  ``$LANES_CONFIG_DIR/lanes.yml``
* ``LANES_REGION``: the AWS region to use when listing EC2 instances. Default:
  ``us-west-2``
* ``LANES_DISABLE_UTF8``: set this to any value to use ASCII for table borders.
  UTF-8 borders are enabled by default.
* ``LANES_TAG_LANE``: the EC2 instance tag to use for determining which lane an
  instance belongs to. Default: ``Lane``
* ``LANES_TAG_NAME``: the EC2 instance tag to use for determining an instance's
  name. Default: ``Name``

## Usage

### Initializing New Systems

The first time you use ``lanes``, you will need to create the initial
configuration files. For this, you can use the ``lanes init`` command.

```bash
# initialize a lanes and create a sample lanes profile
$ lanes init

# initialize a lanes and but do not create a sample lanes profile
$ lanes init --no-profile

# initialize a lanes, overwriting any existing lanes configuration (the
# "default" lanes profile will NOT be overwritten if it exists)
$ lanes init --force
```

Alternatively, you may copy the ``$HOME/.lanes/`` directory from another system
where you have previously configured ``lanes``.

### Creating New Lane Profiles

``lanes`` includes a helper to create fresh lane profiles:

```bash
# create a new profile, prompting for the profile name and AWS credentials
$ lanes init profile

# create a new profile named "foo", prompting only for the AWS credentials
$ lanes init profile foo

# create a new profile named "foo" with "ABCD" as the AWS Access Key ID,
# prompting only for the AWS Secret Access Key
$ lanes init profile foo ABCD
```

Profiles created with this command will include examples for how to configure
individual lanes.

### Selecting Lane Profiles

When executing ``lanes``, the desired profile is determined first by the
``LANES_PROFILE`` environment variable. If this is not set, the profile
configured in ``$HOME/.lanes/lanes.yml`` will be used.

If you wish to quickly change your default profile, you may use ``lanes switch
[new profile name]``.

Examples:

```bash
# override current profile for a single invocation
$ LANES_PROFILE=demo lanes ls

# override current profile for the rest of the terminal session
$ export LANES_PROFILE=demo
$ lanes ls

# set the default profile to $HOME/.lanes/home-profile.yml
$ lanes switch home-profile
```

### Listing EC2 Instances

Examples:

```bash
# list all instances for the current profile
$ lanes list
$ lanes ls

# list all instances in the "prod" lane for the current profile
$ lanes list prod
$ lanes ls prod
```

### SSH Into Instance

Examples:

```bash
# list all instances, prompting for the instance to connect to
$ lanes ssh

# list all instances in the "prod" lane, prompting for the instance to connect to
$ lanes ssh prod
```

### Execute Command On All Lane Instances

Examples:

```bash
# list all instances in the "prod" lane, confirming before executing the
# specified command on each instance
$ lanes sh prod 'ls -l'

# list all instances in the "prod" lane, executing the specified command on
# each instance without confirmation
$ lanes sh prod --confirm 'ls -l'
```

### Push Files to All Lane Instances

Examples:

```bash
# list all instances in the "dev" lane, confirming before copying localfile.txt
# to /tmp/localfile.txt on all instances
$ lanes file push dev localfile.txt /tmp/

# list all instances in the "dev" lane, confirming before copying localfile.txt
# and magic.log to /tmp/ on all instances
$ lanes file push dev localfile.txt magic.log /tmp/

# list all instances in the "dev" lane, copying localfile.txt and magic.log to
# /tmp/ on all instances without confirmation
$ lanes file push dev --confirm localfile.txt magic.log /tmp/
```

## Contributing

To build and install ``lanes`` locally, you will need to have [Go
1.8](https://golang.org/dl/) or newer, as well as [Glide](http://glide.sh) to
manage the build dependencies.

Clone the repository:

```shell
$ git clone https://github.com/codekoala/go-aws-lanes.git
```

Install dependencies:

```shell
$ glide install
```

Build the binary:

```shell
# for Linux systems
$ make linux

# for OSX systems
$ make osx

# for both Linux and OSX
$ make
```

The resulting binaries will appear as ``./bin/lanes_$GOOS_$GOARCH``.

If you just want to run the tests:

```shell
$ make test
```

If you wish to contribute changes to the project, please fork the repository,
make the changes in your fork, and submit a pull request.

## Credits

This project is heavily based on https://github.com/Lemniscate/aws-lanes. The
main reason for building this version was to ease the burden of installing the
utility on different platforms.
