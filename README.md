# go-aws-lanes

This command line tool is designed to help people interact with different sets
of AWS profiles and EC2 instances. It allows you to easily switch between
multiple sets of AWS credentials and perform the following operations:

* list EC2 instances on the account, optionally filtered by a "Lane" tag.
* quickly SSH into a specific EC2 instance using the correct credentials,
  optionally setting up tunnels to locally access services running on a given
  instance.
* copy files to/from one or more EC2 instances (TODO)
* run commands on one or more EC2 instances (TODO)

## Sample Output

```
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
```

## What Are Lanes?

A lane is basically a logical environment for your EC2 instances. For example,
you could have a lane called "dev" for development servers, one called "uat"
user acceptance testing, and one called "prod" for production servers.

## Credits

This project is heavily based on https://github.com/Lemniscate/aws-lanes. The
main reason for building this version was to ease the burden of installing the
utility on different platforms.