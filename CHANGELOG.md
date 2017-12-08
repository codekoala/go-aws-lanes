## v0.4.0 (2017.12.07)

* #19 - Support for running commands in parallel using `lanes sh`
* #20 - Disable tunnels when using `lanes sh`
* Support for CSV output from `lanes ls` using the `--csv` argument

## v0.3.1 (2017.11.06)

* Automatically enable batch mode when piping `lanes ls` output to another
  program

## v0.3.0 (2017.11.01)

* #11 - Automatically switch to a new profile upon creation
* #13 - Support a default profile for instances with no lane specified
* #14 - EC2 instance state included in server table
* #15 - New `lanes edit [profile]` helper command to quickly edit a lanes profile
* #16 - EC2 instances can be filtered by simple keyword matching
* #17 - Table output can be customized
* #18 - `~` is explicitly expanded to the user's home directory
* #9 - Initial support for some bash completion
* Added batch mode to `lanes list` for easier use in scripts
* Added `lanes profiles` to list all known profiles

## v0.2.1 (2017.10.06)

* Releases are built with Go 1.9
* #7 - Check profile permissions, so AWS credentials aren't too easy to find
* #8 - Automatically initialize new systems
* #12 - Clarify expectation with "Which server?" prompt
* Automatically select server when only one is listed

## v0.2.0 (2017.08.04)

* No longer displaying input when prompting for AWS credentials
* #6 - added helpers to create new configuration and lane profiles
* #1 - EC2 instances are now sorted case-insensitively
* #5 - ``lanes.yml`` includes more configuration options
* #4 - tables of EC2 instances can now be printed either using ASCII or UTF-8 border
* #3 - lots of documentation

## v0.1.1 (2017.08.03)

* Fixed bug on systems that never had lanes profiles
* Support for custom tag names to determine instance name and lane
* Released binaries are compressed with UPX by default

## v0.1.0 (2017.08.03)

* Initial release
