+++
title = "Getting started"
description = "A guide to help you get started on using Enbas."
weight = 2
slug = "getting-started"
template = "project-page.html"
+++

# Getting started

## Summary

In this guide we are going to log into an account on a private GoToSocial server.

Follow along to log into your own account.

## Your configuration file

You can use the `--config` top-level flag to specify the path to your configuration file.
If the `--config` flag is not set Enbas will attempt to calculate the default path to your configuratio file
based on your home configuration directory using Go's [os.UserConfigDir()](https://pkg.go.dev/os#UserConfigDir) function.

If you've set the `XDG_CONFIG_HOME` environment variable, the default path to your configuration file
will be set to `$XDG_CONFIG_HOME/enbas/config.json`.

If `XDG_CONFIG_HOME` is not set, then:

- on Linux the path will be set to `$HOME/.config/enbas/config.json`.
- on Darwin (MacOS) the path will be set to `$HOME/Library/Application Support/enbas/config.json`.
- on Windows the path will be set within the `%AppData%` directory.

## Generate your configuration file

Run the following command to generate your configuration file.

```bash
enbas create config
```

Use the `--config` flag to specify the path to save the configuration to.

```bash
enbas --config ./config/config.json create config
```

Once the initial configuration is created you can edit the file to your preferences. 
The [configuration reference page](@/projects/enbas/03_configuration.md) should help you with this.

For this 'Getting Started' guide you may want to specify your preferred browser in the configuration to allow
Enbas to open the link to your instance's authorisation page.
If you prefer to open the link manually then you can leave it empty.

## Log into your GoToSocial account

Enbas uses the Oauth2 authentication flow to log into your account on GoToSocial. Follow the below steps to log into your account:

1. Run the following command to begin the login process. Use the `--url` specifying the instance that you want to log into.
   Use the `--url` flag to specify the URL of the instance you want to log into.
   Use the `--scope` flag to specify the application scope(s) (e.g. read, write).
    ```bash
    enbas create access --url gts.enbas-demo.private --scope read --scope write
    ```

2. Enbas will send a registration request to your instance and receive a new client ID and secret that it
   needs for authentication.

3. Enbas will then generate a link to the consent form for you to access in your browser and print it to
   your terminal screen along with a message explaining that you need to obtain the `out-of-band` token
   to continue.

   The link will open in a tab in your preferred browser if you've specified it in your configuration,
   otherwise you can manually open it yourself.

   If the browser tab doesn't open for you as expected you can still manually open it yourself.

4. Once you've signed into GoToSocial on your browser,
   you will be informed that Enbas would like to perform actions on your behalf.
   If you're happy with this then click on the `Allow` button.
   ![A screenshot of the consent form](/projects/enbas/consent_form.png "Consent Form")

5. The `out-of-band` token from your instance will be displayed to you in your browser.
   Copy it and return to your terminal.

6. Paste the token into the prompt and press `ENTER`.
   Enbas will then exchange the token for an access token which will be used to authenticate
   to your instance on your behalf.

7. Enbas will then verify the access token, save the credentials to the `credentials.json` file
   in your configuration directory, and inform you that you have successfully logged into your account.

### Example login flow

```
$ enbas create access --url gts.enbas-demo.private --scope read --scope write

You'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with the application's login process.
Your browser may have opened the link to the consent page already. If not, please copy and paste the link below to your browser:

https://gts.enbas-demo.private/oauth/authorize?client_id=019RD0WVA903F773T5F9D9EYHP&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=read+write

Once you have the code please copy and paste it below.
Out-of-band token: ZDRKOTE0NMUTZGVHZC0ZNJVJLWJINTMTMWE1M2UWYWFHOTQY
✔ You have successfully logged as percy@gts.enbas-demo.private.
```

## View your account information

You can verify that you have successfully logged in by viewing your account information.

```bash
enbas show account --my-account
```

### Example

```
$ enbas show account --my-account

Percy Cade (@percy)

ACCOUNT ID:
01629QXYN8X597CZDAH4BTY32R

JOINED ON:
29 Jun 2024

STATS:
Followers: 0
Following: 0
Statuses: 0

BIOGRAPHY:
Hey there, the name's Percy.

I've been a Platform Engineer in the Healthcare industry for 5 years and
counting. I love containerising anything and everything with Docker and
Kubernetes, and often find myself dabbling with Python, Go and Rust.

In my free time I like to cook, blog about FOSS software news and make videos
documenting my travels across the UK.


METADATA:
Pronouns: he/him
Location: Hertfordshire, UK
My website: https://percycade.me.private
My blogs: https://blogs.percycade.me.private
My photos: https://photos.percycade.me.private
My videos: https://videos.percycade.me.private

ACCOUNT URL:
https://gts.enbas-demo.private/@percy
```

Now that you have successfully logged into your account proceed to the [user manual](@/projects/enbas/04_user_manual.md).
