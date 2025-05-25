# Enbas Installation Guide

## Table of Contents

- [Summary](#summary)
- [Build requirements](#build-requirements)
- [Obtain the source code](#obtain-the-source-code)
  - [Clone the repository](#clone-the-repository)
  - [Or download the source code](#or-download-the-source-code)
- [Build](#build)
- [Install](#install)
  - [System-wide installation](#system-wide-installation)
  - [Or local installation](#or-local-installation)
- [Verify](#verify)

## Summary

This project does not produce pre-built binaries in order to give you complete control when building and
installing the application. The installation process is pretty simple, but the guide below provides
step-by-step instructions to help you out.

## Build requirements

- **go:**
A minimum version of Go 1.24.3 is required for installing `enbas`.
Visit https://go.dev/dl/ to download the latest version.
- **mage:**
Mage is a build tool similar to `make`.
The project includes mage targets for building the binary, the man pages and the example configuration file.
Visit https://magefile.org/ for instructions on how to install mage.
- **git _(optional)_:**
Git is used to calculate the binary version and commit reference.
This is optional as these values can be set manually via environment variables.

## Obtain the source code

You can obtain the source code by cloning the repository or downloading the source code.

### Clone the repository

You can clone the repository from Codeberg, GitHub or Code Flow.

- from Codeberg
   ```bash
   git clone https://codeberg.org/dananglin/enbas.git
   cd enbas
   ```

- from GitHub
   ```bash
   git clone https://github.com/dananglin/enbas.git
   cd enbas
   ```

- or from Code Flow
   ```bash
   git clone https://codeflow.dananglin.me.uk/apollo/enbas.git
   cd enbas
   ```

### Or download the source code

You can download the source code from Codeberg, GitHub or Code Flow.

- from Codeberg
   ```bash
   curl -sL https://codeberg.org/dananglin/enbas/archive/main.tar.gz -o enbas-main.tar.gz
   tar xzvf enbas-main.tar.gz
   cd enbas
   ```

- from GitHub
   ```bash
   curl -sL https://github.com/dananglin/enbas/archive/refs/heads/main.tar.gz -o enbas-main.tar.gz
   tar xzvf enbas-main.tar.gz
   cd enbas-main
   ```

- or from Code Flow
   ```bash
   curl -sL https://codeflow.dananglin.me.uk/apollo/enbas/archive/main.tar.gz -o enbas-main.tar.gz
   tar xzvf enbas-main.tar.gz
   cd enbas
   ```

## Build

Build the binary and documentation using the commands below.

```bash
mage build:binary
mage build:documentation
```

If you've obtained the source code by downloading the TAR archive, you'll need to set the environment variables
to specify the version and commit reference.

```bash
# Example
export ENBAS_APP_VERSION="v0.2.0-60-gb937807"
export ENBAS_APP_COMMIT_REF="b93780742018a42375f46548b8cca968bd28a669"
mage build:binary
mage build:documentation
```

Once the build is successfully you should see the `__build` directory containing the binary, man pages and the
example configuration file.

```
$ tree __build/
__build/
├── bin
│   └── enbas
└── share
    ├── doc
    │   └── enbas
    │       └── examples
    │           └── config.json
    └── man
        ├── man1
        │   └── enbas.1
        ├── man5
        │   └── enbas.5
        └── man7
            └── enbas-topics.7

10 directories, 5 files
```

## Install

You can install `enbas` system-wide or locally in your home directory.

### System-wide installation

The below set of commands assumes that you want to install the binary within the `/usr/local` directory.
Don't forget to use `sudo` if installing to `/usr/local` requires escalated privileges.

```bash
# install the binary
install --mode 0755 __build/bin/enbas /usr/local/bin

# install the man pages
install -D --mode 0644 __build/share/man/man1/enbas.1 /usr/local/share/man/man1/enbas.1
install -D --mode 0644 __build/share/man/man5/enbas.5 /usr/local/share/man/man5/enbas.5
install -D --mode 0644 __build/share/man/man7/enbas-topics.7 /usr/local/share/man/man5/enbas-topics.7

# install the example configuration file
install -D --mode 0644 __build/share/doc/enbas/examples/config.json /usr/local/share/doc/enbas/examples/config.json
```

### Or local installation

The below set of commands assumes that you want to install the binary within the `${HOME}/.local` directory.

```bash
# install the binary
install --mode 0755 __build/bin/enbas ${HOME}/.local/bin

# install the man pages
install -D --mode 0644 __build/share/man/man1/enbas.1 ${HOME}/.local/share/man/man1/enbas.1
install -D --mode 0644 __build/share/man/man5/enbas.5 ${HOME}/.local/share/man/man5/enbas.5
install -D --mode 0644 __build/share/man/man7/enbas-topics.7 ${HOME}/share/man/man5/enbas-topics.7

# install the example configuration file
install -D --mode 0644 __build/share/doc/enbas/examples/config.json ${HOME}/.local/share/doc/enbas/examples/config.json
```

## Verify

Run `enbas` from your terminal to verify that the installation was successful.
You should see the usage documentation.

```
$ enbas

SUMMARY:
  enbas - A GoToSocial client for the terminal.

VERSION:
  v0.2.0-60-gb937807

USAGE:
  enbas [top-level-flags] <action> <target> [flags]

AVAILABLE TARGETS:
  access: your access to your GoToSocial instance
  account: a local or remote account
  accounts: one or accounts
  alias: the shortname to a set of arguments
  aliases: the list of your aliases
  blocked-accounts: the accounts that are blocked by you
  bookmarks: the statuses that you've bookmarked
  config: your configuration
  favourites: the statuses that you've favourited (liked)
  follow-request: the account that is requesting to follow you
  follow-requests: the list of accounts that are requesting to follow you
  followers: the accounts who are following the specified account
  followings: the accounts who the specified account is following
  instance: the GoToSocial instance
  list: a single list
  lists: one or more lists
  media: the media attached to the specified status
  media-attachment: a media attachment that you own
  muted-accounts: the accounts that are muted by you
  note: your private note about an account
  notification: a single notification
  notifications: multiple notifications
  server: the server mode
  status: a single status
  tag: a single tag (hashtag)
  tags: multiple tags (hashtags)
  thread: a status thread
  timeline: your timeline
  token: details of an application token
  tokens: a list of your tokens
  usage: the usage documentation
  version: the application's build information
  votes: the votes(s) to the poll in a status

TOP-LEVEL FLAGS:
  --config: the path to your configuration file
  --no-color: disable the ANSI colour output when displaying the text on screen

Use "enbas show help --target <target>" for more information about a target and
its supported actions and flags.
```

You can also view the application's version and build information by running `enbas version`.

```
$ enbas version --full

Enbas

Version:    v0.2.0-60-gb937807
Git commit: b937807
Go version: go1.24.3
Build date: 2025-05-18T10:36:27Z
```
