+++
title = "Installation guide"
description = "A guide for installing Enbas on your machine."
weight = 1
slug = "installation-guide"
template = "project-page.html"
+++
# Installation guide

## Download

Pre-built binaries will soon be available on the release page on both Codeberg and GitHub.

## Build from source

### Build requirements

- **Go:** A minimum version of Go 1.23.0 is required for installing Enbas.
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

```
$ enbas
SUMMARY:
    enbas - A GoToSocial client for the terminal.

VERSION:
    v0.2.0

USAGE:
    enbas [flags]
    enbas [flags] [command]

COMMANDS:
    accept      Accepts a request (e.g. a follow request)
    add         Adds a resource to another resource
    block       Blocks a resource (e.g. an account)
    create      Creates a specific resource
    delete      Deletes a specific resource
    edit        Edit a specific resource
    follow      Follow a resource (e.g. an account)
    init        Creates a new configuration file in the specified configuration directory
    login       Logs into an account on GoToSocial
    mute        Mutes a specific resource (e.g. an account)
    reject      Rejects a request (e.g. a follow request)
    remove      Removes a resource from another resource
    show        Shows details about a specified resource
    switch      Performs a switch operation (e.g. switching between logged in accounts)
    unblock     Unblocks a resource (e.g. an account)
    unfollow    Unfollows a resource (e.g. an account)
    unmute      Umutes a specific resource (e.g. an account)
    version     Prints the application's version and build information
    whoami      Prints the account that you are currently logged into

FLAGS:
    --help
        print the help message
    --config-dir
        The path to your configuration directory
    --no-color
        Set to true to disable ANSI colour output when displaying text on screen

Use "enbas [command] --help" for more information about a command.
```

You can also view the application's version and build information using the `version` command.
The build information is correctly displayed if you've downloaded the binary from Codeberg or GitHub,
or if you've built it with Mage.

```bash
$ enbas version --full
Enbas

Version:    v0.2.0
Git commit: fa58e5b
Go version: go1.23.0
Build date: 2024-08-29T07:24:53Z
```

Once you have completed the installation proceed to the [Getting Started guide](@/projects/enbas/02_getting_started.md).
