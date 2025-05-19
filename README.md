# Enbas

## Table of Contents

- [Overview](#overview)
  - [Motivation](#motivation)
  - [Requirements](#requirements)
  - [Development](#development)
  - [Licensing](#licensing)
  - [Inspiration](#inspiration)
- [Installation](#installation)
- [Getting started](#getting-started)
  - [Create and edit your configuration file](#create-and-edit-your-configuration-file)
  - [Log into your GoToSocial account](#log-into-your-gotosocial-account)
    - [Example login flow](#example-login-flow)
  - [View your account information](#view-your-account-information)
- [Further documentation](#further-documentation)

## Overview

`enbas` is a CLI application that allows you to interact with your
[GoToSocial](https://docs.gotosocial.org/en/latest/) instance from your terminal.
With `enbas` you can perform tasks such as:

- viewing timelines
- viewing, creating, deleting and boosting statuses
- viewing media attachments from statuses using your favourite media player
- viewing, following, and blocking accounts
- creating and voting in polls
- viewing, creating, editing and deleting your lists
- ...and much more

### Motivation

This project was created from the desire to interact with my GoToSocial instance and explore the Fediverse
from the comfort of my own terminal emulator. This application is developed for those who use GoToSocial and
spend most of their time in the terminal. If you like the idea of scrolling through your timelines in your
favourite terminal emulator instead of a GUI or a browser then `enbas` might interest you.

### Requirements

- **Your favourite terminal emulator:**
For best results choose a modern terminal emulator such as
[Alacritty](https://alacritty.org/),
[Kitty](https://sw.kovidgoyal.net/kitty/),
[Foot](https://codeberg.org/dnkl/foot),
[Ghostty](https://ghostty.org/)
or [Wezterm](https://wezterm.org/).
By default `enbas` uses ANSI colours so make sure your terminal is configured to use your favourite colour
scheme.
- **A nerd font:**
A nerd font is required to correctly display icons such as displaying the boost icon when viewing a status.
Make sure you have a nerd font installed and that your favourite terminal emulator is configured to use it.
See https://www.nerdfonts.com/ for more information about nerd fonts.
- **A browser:**
A browser is needed for you to complete the login process and if you want to view certain resources in the
browser such as an account or a status.
- **A video player _(optional)_:**
A video player is required if you want to play videos from media attachments.
- **An image viewer _(optional)_:**
An image viewer is required if you want to view images from media attachments.
- **An audio player _(optional)_:**
An audio player is required if you want to play audio from media attachments.

### Development

This project is actively developed in [Code Flow](https://codeflow.dananglin.me.uk/apollo/enbas) with
the `main` branch mirrored to the following forges:

- [**Codeberg**](https://codeberg.org/dananglin/enbas)
- [**Radicle**](https://app.radicle.xyz/nodes/seed.radicle.garden/rad:zhqv2orTvTh2x2d7kYky9NhctrpK)
- [**GitHub**](https://github.com/dananglin/enbas)

### Licensing

The licensing information associated with each file is specified in the [REUSE.toml](REUSE.toml) file,
but in general:

- All original source code is licensed under GPL-3.0-or-later.
- All documentation is licensed under CC-BY-4.0.

### Inspiration

This project was inspired by the following projects:

* **[madonctl](https://github.com/McKael/madonctl)**: A Mastodon CLI client written in Go.
* **[tut](https://github.com/RasmusLindroth/tut)**: A Mastodon TUI written in Go.
* **[toot](https://pypi.org/project/toot/)**: A Mastodon CLI and TUI written in Python.

## Installation

See [INSTALL.md](./INSTALL.md)

## Getting started

In this guide we are going to log into an account on a private GoToSocial server.

Follow along to log into your own account.

### Create and edit your configuration file

The configuration for `enbas` is stored in a JSON formatted file.
The default path to the configuration file is set to `$XDG_CONFIG_HOME/enbas/config.json`.

If the `XDG_CONFIG_HOME` environment variable is not set then:

- on Linux the path will be set to `$HOME/.config/enbas/config.json`.
- on Darwin (MacOS) the path will be set to `$HOME/Library/Application Support/enbas/config.json`.
- on Windows the path will be set within the `%AppData%` directory.

Alternatively you can use the `--config` top-level flag to specify a custom path to your configuration file.

Run the following command to generate your configuration file.

```bash
# Create a new configuration using the default path.
enbas create config
```

```bash
# Or create a new configuration using a custom path.
enbas --config ./path/to/your/config.json create config
```

You can use the example configuration to help you edit your own.
You may recall that when installing the application you may have also generated and installed the example
configuration system-wide (at `/usr/local/share/doc/enbas/examples/config.json`) or locally within your home
directory (at `${HOME}/.local/share/doc/enbas/examples/config.json`).
You can also find a copy in the source repository [here](./doc/examples/config.json).

For more details about the configuration options run `man 5 enbas`.

For this 'Getting Started' guide you may want to specify your favourite browser in the configuration so that
`enbas` can open the link to your instance's authorisation page during the login flow.
If you prefer to open the link manually then you can leave this field empty.

### Log into your GoToSocial account

Enbas uses GoToSocial's Oauth2 authentication flow to log into your Fediverse account.
Follow the below steps to complete the login process.

1. Run the following command to begin the login process.
   Use the `--url` to specify the URL of the instance that you want to log into.
   Use the `--scope` flag to specify the scope of your access (e.g. read, write).
   You can use the `--scope` flag multiple times to specify multiple scopes.
      ```
      enbas login --url gts.mydomain.example --scope read --scope write
      ```

2. `enbas` will send a registration request to your instance and receive a new client ID and secret that it
   needs for authentication.

3. `enbas` will then generate a link to the consent form for you to access in your browser and print it to
   your terminal screen along with a message explaining that you need to obtain the `out-of-band` token
   to continue.

   The link will open in a tab in your preferred browser if you've specified it in your configuration.
   Alternatively you can manually open it yourself.

   If the browser tab doesn't open for you as expected you can still open the link manually.

4. Once you've signed into GoToSocial in your browser,
   you will be informed that `enbas` would like to perform actions on your behalf
   with the scopes that you've specified earlier.
   If you're happy with this then click on the big **Allow** button.
   ![A screenshot of the consent form](./doc/assets/consent_form.png "A screenshot of the consent form")

5. The `out-of-band` token from your instance will be displayed to you in your browser.
   Copy it and return to your terminal.

6. Paste the token into the prompt and press `ENTER`.
   `enbas` will then exchange the token for an access token which will be used to authenticate
   to your instance on your behalf.

7. Finally, `enbas` will then verify the access token, save the credentials to the credentials file at the
   path specified in your configuration file, and inform you that you have successfully logged into your
   account.

#### Example login flow

```
$ enbas login --url super-cell.gts.enbas.private --scope read --scope write

You'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with the application's login process.
Your browser may have opened the link to the consent page already. If not, please copy and paste the link below to your browser:

https://super-cell.gts.enbas.private/oauth/authorize?client_id=01C5TAJ1GC1HFH45BV3BNRSZ1M&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=read+write

Once you have the code, please copy and paste it below.
Out-of-band token: NJVMZGMWZMUTNJDKZI0ZZJNLLWI3NDATYTNJYWE0MJBLOGI5
✔ You have successfully logged in as victor@super-cell.gts.enbas.private.
```

### View your account information

You can verify that you have successfully logged in by viewing your account information by running `enbas show account --my-account`.

```
$ enbas show account --my-account

Victor (@victor)

ACCOUNT ID:
01XWASN1G5K23ZBYCZHR4KQS3C

JOINED ON:
24 Apr 2025

STATS:
Followers: 0
Following: 0
Statuses: 0

BIOGRAPHY:
Hey there, the name's Victor.

I've been a Platform Engineer in the Healthcare industry for 7 years and
counting. I love containerising anything and everything with #docker and #k8s,
and often find myself dabbling with #python, #golang and #rust.

In my free time I like to cook, blog about FOSS software news and make videos
documenting my travels across the UK.

METADATA:
Pronouns: he/him
Location: Hertfordshire, UK
Website: https://victor.me.private
Blogs: https://blogs.victor.me.private
Photos: https://photos.victor.me.private
Streams: https://streams.victor.me.private

ACCOUNT URL:
https://super-cell.gts.enbas.private/@victor

YOUR PREFERENCES:
Default post language: en
Default post visibility: public
Mark posts as sensitive by default: false
```

## Further documentation

All documentation for `enbas` can be viewed using the `man` command.

- To view the general operation manual, run `man 1 enbas`.
- To view the configuration manual for more details about the configuration options, run `man 5 enbas`.
