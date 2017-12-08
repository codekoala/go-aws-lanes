# go-aws-lanes
[![Travis CI Status](https://travis-ci.org/codekoala/go-aws-lanes.svg?branch=master)](https://travis-ci.org/codekoala/go-aws-lanes)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/codekoala/go-aws-lanes/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/codekoala/go-aws-lanes.svg)](https://github.com/codekoala/go-aws-lanes/releases)
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
    ╭──────────────────────────────────────────────────────────╮
    │                       AWS Servers                        │
    ├─────┬──────┬─────────┬────────────┬─────────┬────────────┤
    │ IDX │ LANE │ SERVER  │ IP ADDRESS │ STATE   │ ID         │
    ├─────┼──────┼─────────┼────────────┼─────────┼────────────┤
    │ 1   │ dev  │ dev-01  │ 1.2.3.4    │ running │ i-12341234 │
    │ 2   │ uat  │ uat-01  │ 1.2.3.5    │ running │ i-12341235 │
    │ 3   │ prod │ prod-01 │ 1.2.3.6    │ running │ i-12341236 │
    │ 4   │ prod │ prod-02 │ 1.2.3.7    │ running │ i-12341237 │
    ╰─────┴──────┴─────────┴────────────┴─────────┴────────────╯

    $ lanes ls dev
    Current profile: foo
    Fetching servers... done
    ╭──────────────────────────────────────────────────────────╮
    │                  AWS Servers                             │
    ├─────┬──────┬─────────┬────────────┬─────────┬────────────┤
    │ IDX │ LANE │ SERVER  │ IP ADDRESS │ STATE   │ ID         │
    ├─────┼──────┼─────────┼────────────┼─────────┼────────────┤
    │ 1   │ dev  │ dev-01  │ 1.2.3.4    │ running │ i-12341234 │
    ╰─────┴──────┴─────────┴────────────┴─────────┴────────────╯

    $ lanes ls prod
    Current profile: foo
    Fetching servers... done
    ╭──────────────────────────────────────────────────────────╮
    │                  AWS Servers                             │
    ├─────┬──────┬─────────┬────────────┬─────────┬────────────┤
    │ IDX │ LANE │ SERVER  │ IP ADDRESS │ STATE   │ ID         │
    ├─────┼──────┼─────────┼────────────┼─────────┼────────────┤
    │ 1   │ prod │ prod-01 │ 1.2.3.6    │ running │ i-12341236 │
    │ 2   │ prod │ prod-02 │ 1.2.3.7    │ running │ i-12341237 │
    ╰─────┴──────┴─────────┴────────────┴─────────┴────────────╯

## What Are Lanes?

A lane is basically a logical environment for your EC2 instances. For example,
you could have a lane called "dev" for development servers, one called "uat"
user acceptance testing, and one called "prod" for production servers.

## Installation

### Manual Installation

1. Download a pre-compiled, released version from [the releases
   page](https://github.com/codekoala/go-aws-lanes/releases)
2. Mark the binary as executable
3. Move the binary to your ``$PATH``

For example:

```shell
$ curl -Lo /tmp/lanes https://github.com/codekoala/go-aws-lanes/releases/download/v0.4.0/lanes_linux_amd64
$ chmod +x /tmp/lanes
$ sudo mv /tmp/lanes /usr/local/bin/lanes
```

To compile from source, please see the instructions in the [contributing section](#contributing).

### Arch Linux

There is a PKGBUILD in the [AUR](https://aur.archlinux.org/packages/aws-lanes/)
to help package and install ``lanes`` for Arch Linux-based distributions.

## Usage

### Initializing New Systems

As of v0.2.1, initial configuration is handled the first time you run any lanes
command. You may also use the ``lanes init`` command.

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
individual lanes. ``lanes`` automatically switches to profiles created with
this command. If you would like to create a new profile without switching to it
immediately, use the ``--no-switch`` or ``-n`` flags:

```bash
# create a new profile named "foo", prompting only for the AWS credentials,
# without automatically switching to the new profile
$ lanes init profile foo --no-switch
```

### Editing Lane Profiles

``lanes`` includes a helper to quickly open the configuration for a specific
profile in your default editor. Your default editor is determined by the
`$EDITOR` environment variable. If this variable is not set, ``lanes`` attempts
to use ``vi``.

```bash
# edit your current profile using your default editor
$ lanes edit

# edit the profile called "foo" using your default editor
$ lanes edit foo
```

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

As of version 0.3.0, the `list`/`ls` command has a `--batch`/`-b` option to
disable table headers and borders for easier use with batch operations. It is
also possible to show specific columns with the `--columns`/`-c` option.
Alternatively, specific columns may be hidden using the `--hide` option.

Using the `list` command in batch mode can be helpful when writing other
scripts to interact with your AWS EC2 instances. For example, here's a
one-liner to produce a roster for `salt-ssh`:

```bash
$ lanes ls -c SSH_IDENTITY,USER,IP,NAME | \
    sed "s,~,$HOME,g" | \
    awk '/\.pem/ { \
print $4":\n \
  host: "$3"\n \
  user: "$2"\n \
  sudo: true\n \
  tty: true\n \
  priv: "$1"\n \
"}' > /etc/salt/roster
```

As of version 0.4.0, the `list`/`ls` command also supports dumping the server
table in CSV format using the `--csv` argument.

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

As of version 0.4.0, `lanes sh` supports running the specified command on
multiple machines in parallel. There are three different options to enable
parallel execution:

* `--parallel` runs the specified command on all instances in the specified
  lane at the same time.
* `--num-parallel/-n N` runs the specified command on up to `N` instances in
  the specified line at the same time.
* `--pparallel N` runs the specified command on up to `N%` of the instances in
  the specified lane at the same time.

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

## Configuration

The configuration for this tool lives in ``$HOME/.lanes/`` by default. There
are two forms of configuration for ``lanes``: the configuration for ``lanes``
itself and configuration for individual lanes in their respective files.

The configuration for ``lanes`` itself lives in ``$HOME/.lanes/lanes.yml`` by
default. Here are the configuration options:

```yaml
profile: default
region: us-west-2
disable_utf8: false
tags:
  name: Name
  lane: Lane
```

* ``profile: default``: this indicates that the "lane profile" should be read
  from ``$HOME/.lanes/default.yml``.
* ``region: us-west-2``: this is the default AWS region to use when querying
  EC2 instances.
* ``disable_utf8: false``: this setting can be used to toggle UTF-8 and ASCII
  mode for table borders.
* ``tags.name: Name``: this indicates that the EC2 instance tag named "Name"
  will be used to determine each instance's name. Change this if you use a
  different tag name in your environment.
* ``tags.lane: Lane``: this indicates that the EC2 instance tag named "Lane"
  will be used to determine each instance's lane. Change this if you use a
  different tag name in your environment.

The configuration for an individual lane lives in ``$HOME/.lanes/[lane profile
name].yml`` by default. Here are the configuration options:

```yaml
aws_access_key_id: ASDF
aws_secret_access_key: FDSA
region: us-east-1
ssh:
  mods:
    dev:
      identity: ~/.ssh/id_rsa_dev
      tunnels:
        - 8080:127.0.0.1:80
        - 3306:127.0.0.1:3306
    uat:
      identity: ~/.ssh/id_rsa_uat
      tunnel: 8080:127.0.0.1:80
    prod:
      identity: ~/.ssh/id_rsa_prod
```

* ``aws_access_key_id``: the AWS access key ID for the lane profile.
* ``aws_secret_access_key``: the AWS secret access key for the lane profile.
* ``region``: the default region for this lane profile. If not specified, the
  region will be determined by the global configuration for ``lanes`` (see
  above).
* ``ssh.mods.[lane name].user``: the username to use when SSH'ing into an EC2
  instance in the specified lane.
* ``ssh.mods.[lane name].identity``: the private key to use when SSH'ing into
  instances in the specified lane.
* ``ssh.mods.[lane name].tunnel``: a single tunnel to setup when SSH'ing to a
  specific EC2 instance in the specified lane.
* ``ssh.mods.[lane name].tunnels``: a list of tunnels to setup when SSH'ing to
  a specific EC2 instance in the specified lane.

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
