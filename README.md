# go-aws-lanes
[![Travis CI Status](https://travis-ci.org/codekoala/go-aws-lanes.svg?branch=master)](https://travis-ci.org/codekoala/go-aws-lanes)
[![License BSD3](https://img.shields.io/badge/license-BSD3-blue.svg)](https://raw.githubusercontent.com/codekoala/go-aws-lanes/master/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/codekoala/go-aws-lanes/total.svg)](https://github.com/codekoala/go-aws-lanes/releases)

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
    ╭─────────────────────────────────────────────────────────────────────────────────────────────╮
    │                                       AWS Servers                                           │
    ├─────┬──────────────────┬──────────────────────────────────────────────┬─────────────────────┤
    │ IDX │ LANE             │ SERVER                      │ IP ADDRESS     │ ID                  │
    ├─────┼──────────────────┼─────────────────────────────┼────────────────┼─────────────────────┤
    │ 1   │ dev              │ dev-01                      │ 1.2.3.4        │ i-12341234          │
    │ 2   │ uat              │ uat-01                      │ 1.2.3.5        │ i-12341235          │
    │ 3   │ prod             │ prod-01                     │ 1.2.3.6        │ i-12341236          │
    │ 4   │ prod             │ prod-02                     │ 1.2.3.7        │ i-12341237          │
    ╰─────┴──────────────────┴──────────────────────────────────────────────┴─────────────────────╯

    $ lanes ls dev
    Current profile: foo
    Fetching servers... done
    ╭─────────────────────────────────────────────────────────────────────────────────────────────╮
    │                                       AWS Servers                                           │
    ├─────┬──────────────────┬──────────────────────────────────────────────┬─────────────────────┤
    │ IDX │ LANE             │ SERVER                      │ IP ADDRESS     │ ID                  │
    ├─────┼──────────────────┼─────────────────────────────┼────────────────┼─────────────────────┤
    │ 1   │ dev              │ dev-01                      │ 1.2.3.4        │ i-12341234          │
    ╰─────┴──────────────────┴──────────────────────────────────────────────┴─────────────────────╯

    $ lanes ls prod
    Current profile: foo
    Fetching servers... done
    ╭─────────────────────────────────────────────────────────────────────────────────────────────╮
    │                                       AWS Servers                                           │
    ├─────┬──────────────────┬──────────────────────────────────────────────┬─────────────────────┤
    │ IDX │ LANE             │ SERVER                      │ IP ADDRESS     │ ID                  │
    ├─────┼──────────────────┼─────────────────────────────┼────────────────┼─────────────────────┤
    │ 1   │ prod             │ prod-01                     │ 1.2.3.6        │ i-12341236          │
    │ 2   │ prod             │ prod-02                     │ 1.2.3.7        │ i-12341237          │
    ╰─────┴──────────────────┴──────────────────────────────────────────────┴─────────────────────╯

## What Are Lanes?

A lane is basically a logical environment for your EC2 instances. For example,
you could have a lane called "dev" for development servers, one called "uat"
user acceptance testing, and one called "prod" for production servers.

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

### Environment Variables

``lanes`` supports a handful of environment variables to quickly change
behavior:

* ``LANES_CONFIG_DIR``: the directory where all configuration is expected to
  reside. Default: ``$HOME/.lanes/``
* ``LANES_CONFIG``: the configuration file to use for lanes. Default:
  ``$LANES_CONFIG_DIR/lanes.yml``
* ``LANES_REGION``: the AWS region to use when listing EC2 instances. Default:
  ``us-west-2``

## Usage

When executing ``lanes``, the desired profile is determined first by the
``LANES_PROFILE`` environment variable. If this is not set, the profile
configured in ``$HOME/.lanes/lanes.yml`` will be used.

## Credits

This project is heavily based on https://github.com/Lemniscate/aws-lanes. The
main reason for building this version was to ease the burden of installing the
utility on different platforms.
