# Installation Guide

## Download

Pre-built binaries will soon be available on the release page on both Codeberg and GitHub.

## Build from source

### Build requirements

- **Go:** A minimum version of Go 1.23.0 is required for installing spruce.
  Please go [here](https://go.dev/dl/) to download the latest version.

- **Mage (Optional):** The project includes mage targets for building and installing the binary. The main
  advantage of using mage over `go install` is that the build information is baked into the binary during
  compilation. You can visit the [Mage website](https://magefile.org/) for instructions on how to install Mage.

### Install with mage

You can install Enbas with Mage using the following commands:

```bash
git clone https://github.com/dananglin/enbas.git
mage install
```

By default Mage will attempt to install Enbas to `/usr/local/bin/enbas` which will most likely fail as you won't
the permission to write to `/usr/local/bin/`. You will need to either run `sudo mage install`, or you can
(preferably) change the install prefix to a directory that you have permission to write to using
the `ENBAS_INSTALL_PREFIX` environment variable. For example:

```bash
ENBAS_INSTALL_PREFIX=${HOME}/.local mage install
```

This will install Enbas to `~/.local/bin/enbas`.

The table below shows all the environment variables you can use when building with Mage.

| Environment Variable     | Description                                                                  |
|--------------------------|------------------------------------------------------------------------------|
|`ENBAS_INSTALL_PREFIX`    | Set this to your preferred the installation prefix (default: `/usr/local`).  |
|`ENBAS_BUILD_REBUILD_ALL` | Set this to `1` to rebuild all packages even if they are already up-to-date. |
|`ENBAS_BUILD_VERBOSE`     | Set this to `1` to enable verbose logging when building the binary.          |

### Install with go

If your `GOBIN` directory is included in your `PATH` then you can install Enbas with Go.

```bash
git clone https://github.com/dananglin/enbas.git
cd enbas
go install ./cmd/enbas
```

## Verify the installation

Type `enbas` from your terminal to verify that the installation was successful. You should see the help documentation.

```bash
$ enbas
SUMMARY:
    enbas - A GoToSocial client for the terminal.

VERSION:
    v0.1.0

USAGE:
    enbas [flags]
    enbas [flags] [command]

COMMANDS:
    accept      Accept a request (e.g. a follow request)
    add         Add a resource to another resource
    block       Block a resource (e.g. an account)
    create      Create a specific resource
    delete      Delete a specific resource
    edit        Edit a specific resource
    follow      Follow a resource (e.g. an account)
    init        Create a new configuration file in the specified configuration directory
    login       Login to an account on GoToSocial
    mute        Mute a resource (e.g. an account)
    reject      Reject a request (e.g. a follow request)
    remove      Remove a resource from another resource
    show        Print details about a specified resource
    switch      Perform a switch operation (e.g. switch logged in accounts)
    unblock     Unblock a resource (e.g. an account)
    unfollow    Unfollow a resource (e.g. an account)
    unmute      Unmute a resource (e.g. an account)
    version     Print the application's version and build information
    whoami      Print the account that you are currently logged in to

FLAGS:
    --help
        print the help message
    --config-dir
        Specify your config directory
    --no-color
        Disable ANSI colour output when displaying text on screen

Use "enbas [command] --help" for more information about a command.
```

You can also view the application's version and build information using the `version` command.
The build information is correctly displayed if you've downloaded the binary from Codeberg or GitHub,
or if you've built it with Mage.

```bash
Enbas

Version:    v0.1.0
Git commit: c8892a6
Go version: go1.22.4
Build date: 2024-06-28T21:57:37Z
```

Once you have completed the installation proceed to the [Getting Started guide](./getting_started.md).
