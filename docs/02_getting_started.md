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

## Your configuration directory

You can use the `--config-dir` global flag to specify the path to your configuration directory.
Alternatively Enbas tries to set the directory based on your home configuration directory using Go's [os.UserConfigDir()](https://pkg.go.dev/os#UserConfigDir) function.

If you've set the `XDG_CONFIG_HOME` environment variable, the configuration directory will be set to `$XDG_CONFIG_HOME/enbas`.

If this is not set, then:

- on Linux the configuration directory will be set to `$HOME/.config/enbas`.
- on Darwin (MacOS) the configuration directory will be set to `$HOME/Library/Application Support/enbas`.
- on Windows the configuration directory will be set within the `%AppData%` directory.

## Generate your configuration file

Run the `init` command to generate your configuration file.

```bash
enbas init
```

Use the `--config-dir` flag if you want to generate it in a specific directory

```bash
enbas --config-dir ./config init
```

You should now see a file called `config.json` in your configuration directory.
Feel free to edit the file to your preferences. 
The [configuration reference page](@/projects/enbas/03_configuration.md) should help you with this.

For this 'Getting Started' guide you may want to specify your preferred browser in the configuration to allow
Enbas to open the link to your instance's authorisation page.
If you prefer to open the link manually then you can leave it blank.

## Log into your GoToSocial account

Enbas uses the Oauth2 authentication flow to log into your account on GoToSocial.

<details>
    <summary style="color: Orange; font-weight: bold"><span>&#9888;</span> Warning</summary>
    <p>
        As of writing GoToSocial does not currently support scoped authorization tokens so even if
        we request read-only tokens, the application will be able to perform any actions within the
        limitations of your account (including admin actions if you are an admin).
        You can read more about this <a href="https://docs.gotosocial.org/en/latest/api/authentication/">here</a>.
    </p>
</details>

Follow the below steps to log into your account:

1. Run the `login` command specifying the instance that you want to log into.
    ```bash
    enbas login --instance gts.enbas-demo.private
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
$ enbas login --instance gts.enbas-demo.private

You'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with the application's login process.
Your browser may have opened the link to the consent page already. If not, please copy and paste the link below to your browser:

https://gts.enbas-demo.private/oauth/authorize?client_id=019RD0WVA903F773T5F9D9EYHP&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code

Once you have the code please copy and paste it below.
Out-of-band token: ZDRKOTE0NMUTZGVHZC0ZNJVJLWJINTMTMWE1M2UWYWFHOTQY
✔ You have successfully logged as percy@gts.enbas-demo.private.
```

## View your account information

You can verify that you have successfully logged in by viewing your account information.

```bash
enbas show --type account --my-account
```

### Example

```
$ enbas show --type account --my-account

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
