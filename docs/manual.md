# User Manual

## Table of Contents

- [Global flags](#global-flags)
- [Version](#version)
  - [Print the application version](#print-the-application-version)
- [Init](#init)
- [Authentication](#authentication)
  - [Logging into an account](#logging-into-an-account)
  - [Switch between accounts](#switch-between-accounts)
  - [See the account that you are currently logged in as](#see-the-account-that-you-are-currently-logged-in-as)
- [Accounts](#accounts)
  - [View account information](#view-account-information)
  - [Follow an account](#follow-an-account)
  - [Unfollow an account](#unfollow-an-account)
  - [Show an account's followers](#show-an-accounts-followers)
  - [Show account's followings](#show-accounts-followings)
  - [Block an account](#block-an-account)
  - [Unblock an account](#unblock-an-account)
  - [View blocked accounts](#view-blocked-accounts)
  - [Mute an account](#mute-an-account)
  - [Unmute an account](#unmute-an-account)
  - [View muted accounts](#view-muted-accounts)
  - [Add a private note to an account](#add-a-private-note-to-an-account)
  - [Remove the private note from an account](#remove-the-private-note-from-an-account)
- [Follow requests](#follow-requests)
  - [View your follow requests](#view-your-follow-requests)
  - [Accept a follow request](#accept-a-follow-request)
  - [Reject a follow request](#reject-a-follow-request)
- [Media Attachments](#media-attachments)
  - [Create a media attachment](#create-a-media-attachment)
  - [Edit a media attachment](#edit-a-media-attachment)
  - [View a media attachment](#view-a-media-attachment)
- [Statuses](#statuses)
  - [View a status](#view-a-status)
  - [Create a status](#create-a-status)
  - [Delete a status](#delete-a-status)
  - [Boost (Repost) a status](#boost-repost-a-status)
  - [Un-boost (Un-repost) a status](#un-boost-un-repost-a-status)
  - [Like a status](#like-a-status)
  - [Unlike a status](#unlike-a-status)
  - [View a list of statuses that you've liked](#view-a-list-of-statuses-that-youve-liked)
  - [Mute a status](#mute-a-status)
  - [Unmute a status](#unmute-a-status)
  - [Vote in a poll within a status](#vote-in-a-poll-within-a-status)
- [Polls](#polls)
  - [Create a poll](#create-a-poll)
  - [View a poll](#view-a-poll)
  - [Vote in a poll](#vote-in-a-poll)
- [Lists](#lists)
  - [Create a list](#create-a-list)
  - [View a list of your lists](#view-a-list-of-your-lists)
  - [View a specific list](#view-a-specific-list)
  - [Edit a list](#edit-a-list)
  - [Delete a list](#delete-a-list)
  - [Add accounts to a list](#add-accounts-to-a-list)
  - [Remove accounts from a list](#remove-accounts-from-a-list)
- [Timelines](#timelines)
  - [View a timeline](#view-a-timeline)
- [Media](#media)
  - [View media from a status](#view-media-from-a-status)
- [Bookmarks](#bookmarks)
  - [View your bookmarks](#view-your-bookmarks)
  - [Add a status to your bookmarks](#add-a-status-to-your-bookmarks)
  - [Remove a status from your bookmarks](#remove-a-status-from-your-bookmarks)
- [Notifications](#notifications)

## Global flags

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `config-dir` | string | false | The configuration directory. | |
| `no-color` | boolean | false | Disables ANSI colour output when displaying text on screen<br>You can also set `NO_COLOR` to any value for the same effect. | false |

## Version

### Print the application version

View the application's version and build information

```
enbas version --full
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `full` | boolean | false | Prints the full build information. | false |

## Init

Initialises the app by creating a configuration file in the configuration directory.
If you want to use a specific directory then use the global `--config-dir` flag.

```
enbas init
```

## Authentication

### Logging into an account

Log into your GoToSocial account. You can run this multiple times to log into multiple accounts.

```
enbas login --instance gts.enbas-demo.private
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `instance` | string | true | The instance that you want to log into. | |

### Switch between accounts

Switch between your logged in accounts.

```
enbas switch --to account --account-name vincent@gts.enbas-demo.private
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `to` | string | true | The resource you want to switch to. In this case you want `account`. | |
| `account-name` | string | true | The name of the account you want to switch to. | |

### See the account that you are currently logged in as

```
enbas whoami
```

## Accounts

### View account information

- View information from your own account
   ```
   enbas show --type account --my-account
   ```
- View information from a local or remote account.
   ```
   enbas show --type account --account-name @name@example.social
   ```
- View an account and show the public statuses that it has created.
   ```
   enbas show --type account --account-name @name@example.social --show-statuses --only-public
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view. Here this should be `account`. | |
| `my-account` | boolean | true | Set to `true` to view your own account. | |
| `show-preferences` | boolean | false | Show your posting preferences. Only applicable with the `my-account` flag. | false |
| `account-name` | string | false | The name of the account to view. This is not required with the `my-account` flag. | |
| `skip-relationship` | boolean | false | Set to `true` to skip viewing your relationship to the account you are viewing (including the private note if you've created one). | false |

Additional flags for viewing an account's statuses.

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `show-statuses` | bool | false | Set to `true` to view the statuses created from this account. | false |
| `limit` | integer | false | The maximum amount of statuses to show from this account. | 20 |
| `exclude-replies` | bool | false | Set to `true` to exclude replies. | false |
| `exclude-boosts` | bool | false | Set to `true` to exclude boosts. | false |
| `only-pinned` | bool | false | Set to `true` to view only pinned statuses. | false |
| `only-media` | bool | false | Set to `true` to view only statuses with media attachments. | false |
| `only-public` | bool | false | Set to `true` to view only public statuses. | false |

### Follow an account

Sends a follow request to the account you want to follow.

```
enbas follow --type account --account-name @name@example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to follow. Here this should be `account`. | |
| `account-name` | string | true | The name of the account to follow. | |
| `show-reposts` | boolean | false | Show reposts from the account you want to follow. | true |
| `notify` | boolean | false | Get notifications when the account you want to follow posts a status. | false |

### Unfollow an account

Unfollows the account that you are currently following.
If you have a follow request pending for the account in question,
performing an unfollow action will remove said follow request.

```
enbas unfollow --type account --account-name @name@example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to unfollow. Here this should be `account`. | |
| `account-name` | string | true | The name of the account to unfollow. | |

### Show an account's followers

- View followers of your own account.
   ```
   enbas show --type followers --from account --my-account
   ```
- View followers of another account.
   ```
   enbas show --type followers --from account --account-name @name@example.social
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view. Here this should be `followers`. | |
| `from` | string | true | The resource you want to view followers from.<br>Here this should be `account`. | |
| `my-account` | boolean | false | Set to `true` to view followers from your own account.<br>This takes precendence over `account-name`.| false |
| `account-name` | string | true | The name of the account to get the followers from. | |

### Show account's followings

- View the accounts that you are following.
   ```
   enbas show --type following --from account --my-account
   ```
- View the accounts that another account is following.
   ```
   enbas show --type following --from account --account-name @name@example.social
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view. Here this should be `following`. | |
| `from` | string | true | The resource you want to view the followings from.<br>Here this should be `account`. | |
| `my-account` | boolean | false | Set to `true` to view the list from your own account.<br>This takes precendence over `account-name`.| false |
| `account-name` | string | true | The name of the account to get the list from. | |

### Block an account

```
enbas block --type account --account-name @name@example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to block. Here this should be `account`. | |
| `account-name` | string | true | The name of the account to block. | |

### Unblock an account

```
enbas unblock --type account --account-name @name@example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to unblock. Here this should be `account`. | |
| `account-name` | string | true | The name of the account to unblock. | |

### View blocked accounts

Prints a list of accounts that you are currently blocking.

```
enbas show --type blocked
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `blocked` for blocked accounts. | |
| `limit` | integer | false | The maximum number of accounts to list. | 20 |

### Mute an account

```
enbas mute --type account --account-name @name@example.social --mute-notifications --mute-duration="1h"
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to mute.<br>Here this should be `account`. | |
| `account-name` | string | true | The name of the account to mute. | |
| `mute-notifications` | boolean | false | Set to `true` to mute notifications as well as statuses. | false |
| `mute-duration` | string | false | Specify how long the account should be muted for.<br>Set to `0s` to mute indefinitely | 0s (indefinitely). |

### Unmute an account

```
enbas unmute --type account --account-name @name@example.social 
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to unmute.<br>Here this should be `account`. | |
| `account-name` | string | true | The name of the account to unmute. | |

### View muted accounts

Prints a list of accounts that you have muted.

``` 
enbas show --type muted-accounts
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `muted-accounts`. | |
| `limit` | integer | false | The maximum number of accounts to print. | 20 |

### Add a private note to an account

Adds a private note to an account. Private notes can only be viewed by you.

```
enbas add --type note --to account --account-name @name@example.social --content "This person is awesome."
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `note`. | |
| `to` | string | true | The resource you want to add the note to.<br>Here this should be `account`. | |
| `account-name` | string | true | The name of the account that you want to add the note to. | |
| `content` | string | true | The content of the note. | |

### Remove the private note from an account

```
enbas remove --type note --from account --account-name @name@example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to remove.<br>Here this should be `note`. | |
| `from` | string | true | The resource you want to remove the note to.<br>Here this should be `account`. | |
| `account-name` | string | true | The name of the account that you want to remove the note from. | |

## Follow requests

### View your follow requests

Prints a list of accounts that are requesting to follow you.

```
enbas show --type follow-request
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `follow-request`. | |
| `limit` | integer | false | The maximum number of accounts to print. | 20 |

### Accept a follow request

Accepts the request from the account that wants to follow you.

```
enbas accept --type follow-request --account-name @person.example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to accept.<br>Here this should be `follow-request`. | |
| `account-name` | string | true | The name of the account that you want to accept. | |

### Reject a follow request

Rejects the request from the account that wants to follow you.

```
enbas reject --type follow-request --account-name @person.example.social
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to accept.<br>Here this should be `follow-request`. | |
| `account-name` | string | true | The name of the account that you want to reject. | |

## Media Attachments

### Create a media attachment

Uploads media from a file to the instance and creates a media attachment.
You can write a description of the media in a text file and specify the path with the `media-description` flag (see the examples below).

- Create a media attachment with a simple description and a focus of x=-0.1, y=0.5
   ```
   enbas create --type media-attachment \
       --media-file picture.png \
       --media-description "A picture of an old, slanted wooden bench in front of the woods." \
       --media-focus "-0.1,0.5"
   ```
- Create a media attachment using a description written in the `description.txt` text file.
   ```
   enbas create --type media-attachment \
       --media-file picture.png \
       --media-description file@description.txt
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to create.<br>Here this should be `media-attachment`. | |
| `media-file` | string | true | The path to the media file. | |
| `media-description` | string | false | The description of the media attachment which will be used as the media's alt-text.<br>To use a description from a text file, use the `flag@` prefix followed by the path to the file (e.g. `file@description.txt`)| |
| `media-focus` | string | false | The media's focus values. This should be in the form of two comma-separated numbers between -1 and 1 (e.g. 0.25,-0.34) | |

### Edit a media attachment

Edits the description and/or the focus of a media attachment that you own.

```
enbas edit --type media-attachment \
    --attachment-id 01J5B9A8WFK59W11MS6AHPYWBR \
    --media-description "An updated description of a picture."
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to edit.<br>Here this should be `media-attachment`. | |
| `media-description` | string | false | The description of the media attachment to edit.<br>To use a description from a text file, use the `flag@` prefix followed by the path to the file (e.g. `file@description.txt`)| |
| `media-focus` | string | false | The media's focus values. This should be in the form of two comma-separated numbers between -1 and 1 (e.g. 0.25,-0.34) | |

### View a media attachment

Prints information about a given media attachment that you own.
You can only see information about the media attachment that you own.

```
enbas show --type media-attachment --attachment-id 01J0N0RQSJ7CFGKHA30F7GBQXT
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `media`. | |
| `attachment-id` | string | true | The ID of the media attachment to view. | |

## Statuses

### View a status

Prints information of a status on screen.
If the `--browser` flag is used, the link to the status is opened instead.
To enable browser support you must specify the browser in your configuration.

See the [configuration reference page](configuration.md#integration) on how to set up integration with
your browser if you have not done so already.

```
enbas show --type status --status-id 01J1Z9PT0243JT9QNQ5W96Z8CA
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status that you want to view. | |
| `browser` | boolean | false | Set to true to open the link to the status in your browser. | false |

### Create a status

- Create a simple status that is publicly visible.
   ```
   enbas create --type status --content-type plain --visibility public --content "Hello, Fediverse!"
   ```
- Create a private status from a file.
   ```
   enbas create --type status --content-type markdown --visibility private --from-file status.md
   ```
- Reply to another status.
  ```
  enbas create --type status --in-reply-to 01J2A86E3M7WWH37H1QENT7CSH --content "@bernie thanks for this! Looking forward to trying this out."
  ```
- Create a status with a poll
   ```
   enbas create \
       --type status \
       --content-type plain \
       --visibility public \
       --content "The age-old question: which text editor do you prefer?" \
       --add-poll \
       --poll-allows-multiple-choices=false \
       --poll-expires-in 168h \
       --poll-option "emacs" \
       --poll-option "vim/neovim" \
       --poll-option "nano" \
       --poll-option "other (please comment)"
   ```
   ![A screenshot of a status with a poll](../assets/images/created_poll.png "A status with a poll")
- Create a status with a media attachment that you have created.
   ```
   enbas create \
       --type status \
       --attachment-id 01J5BDHYJ7MWMMG76FP49H7SWD \
       --content "I went out for a walk in the woods and found this interesting looking wooden bench."
   ```
- Upload and attach 4 media files to a new status. You must set the same number of `media-description` and `media-focus` flags **must** as the `media-file` flags.
  The first `media-description` and `media-focus` flags correspond to the first `media-file` flag and so on.
   ```
   enbas create --type status --visibility public \
       --content "This post has a picture of a cat, a dog, a bee and a bird." \
       --media-file cat.jpg   --media-description file@cat.txt  --media-focus "0,0" \
       --media-file dog.jpg   --media-description file@dog.txt  --media-focus "-0.1,0.25" \
       --media-file bee.jpg   --media-description file@bee.txt  --media-focus "1,1" \
       --media-file bird.webp --media-description file@bird.txt --media-focus "0,0"
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to create.<br>Here this should be `status`. | |
| `attachment-id` | string | false | The ID of the media attachment to attach to the status.<br>Use this flag multiple times to attach multiple media. |
| `content` | string | false | The content of the status.<br>This flag takes precedence over `from-file`.| |
| `content-type` | string | false | The format that the content is created in.<br>Valid values are `plain` and `markdown`. | plain |
| `enable-reposts` | boolean | false | The status can be reposted (boosted) by others. | true |
| `enable-federation` | boolean | false | The status can be federated beyond the local timelines. | true |
| `enable-likes` | boolean | false | The status can be liked (favourtied). | true |
| `enable-replies` | boolean | false | The status can be replied to. | true |
| `from-file` | string | false | The path to the file where to read the contents of the status from. | |
| `in-reply-to` | string | false | The ID of the status that you want to reply to. | |
| `language` | string | false | The ISO 639 language code that the status is written in.<br>If this is not specified then the default language from your posting preferences will be used. | |
| `media-file` | string | false | The path to the media file.<br>Use this flag multiple times to upload multiple media files. | |
| `media-description` | string | false | The description of the media attachment which will be used as the media's alt-text.<br>To use a description from a text file, use the `flag@` prefix followed by the path to the file (e.g. `file@description.txt`)<br>Use this flag multiple times to set multiple descriptions.| |
| `media-focus` | string | false | The media's focus values. This should be in the form of two comma-separated numbers between -1 and 1 (e.g. 0.25,-0.34).<br>Use this flag multiple times to set multiple focus values. | |
| `sensitive` | string | false | The status should be marked as sensitive.<br>If this is not specified then the default sensitivity from your posting preferences will be used. | |
| `spoiler-text` | string | false | The text to display as the status' warning or subject. | |
| `visibility` | string | false | The visibility of the status.<br>Valid values are `public`, `private`, `unlisted`, `mutuals_only` and `direct`.<br>If this is not specified then the default visibility from your posting preferences will be used. | |

Additional flags for polls.

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `add-poll` | boolean | false | Set to `true` to add a poll to the status. | false |
| `poll-allows-multiple-choices` | boolean | false | Set to `true` to allow users to make multiple choices. | false |
| `poll-hides-vote-counts` | boolean | false | Set to `true` to hide the vote count until the poll is closed. | false |
| `poll-option` | string | true | An option in the poll. Use this flag multiple times to set multiple options. | |
| `poll-expires-in` | string | false | The duration in which the poll is open for. | |

### Delete a status

_Not yet supported_

### Boost (Repost) a status

To boost a status, simply add a `boost` to it.

```
enbas add --type boost --to status --status-id 01J17FH1KD9CN6J9Q01011NE0D
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `boost`. | |
| `to` | string | true | The resource you want to add the boost to.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status that you want to boost. | |

### Un-boost (Un-repost) a status

To un-boost a status that you've boosted, simply remove the `boost` from it.

```
enbas remove --type boost --from status --status-id 01J17FH1KD9CN6J9Q01011NE0D
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `boost`. | |
| `from` | string | true | The resource you want to remove the boost from.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status that you want to un-boost. | |

### Like a status

To like (favourite) a status, simply add a `like` or a `star` to it.

```
enbas add --type star --to status --status-id 01J17FH1KD9CN6J9Q01011NE0D
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should either be `like` or `star`. | |
| `to` | string | true | The resource you want to like.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status that you want to like. | |

### Unlike a status

To unlike (un-favourite) a status that you've previously liked, simply remove the `like` or `star` from it.

```
enbas remove --type star --from status --status-id 01J17FH1KD9CN6J9Q01011NE0D
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should either be `like` or `star`. | |
| `from` | string | true | The resource you want to remove the like from.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status that you want to remove the like from. | |

### View a list of statuses that you've liked

Prints the list of statuses that you've liked.

```
enbas show --type liked
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should either be `liked` or `starred`. | |
| `limit` | integer | false | The maximum number of statuses to print. | 20 |

### Mute a status

_Not yet supported_

### Unmute a status

_Not yet supported_

### Vote in a poll within a status

Adds your vote(s) to a poll within a status.

```
enbas add --type vote --to status --status-id 01J55XVV2MM6MKQ7QHFBAVAE8R --vote 3
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `vote`. | |
| `to` | string | true | The resource you want to add the vote to.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the poll you want to add the votes to. | |
| `vote` | int | true | The ID of the option that you want to vote for.<br>You can use this flag multiple times to vote for more than one option if the poll allows multiple choices. | |

## Polls

### Create a poll

See [Create a status](#create-a-status).

### View a poll

You can view a poll within a [status](#view-a-status) or within a [timeline](#view-a-timeline).

### Vote in a poll

See [Vote in a poll within a status](#vote-in-a-poll-within-a-status)

## Lists

### Create a list

```
enbas create --type list --list-title "My Favourite People" --list-replies-policy list
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to create.<br>Here this should be `list`. | |
| `list-title` | string | true | The title of the list that you want to create. | |
| `list-replies-policy` | string | false | The policy of the replies for this list.<br>Valid values are `followed`, `list` and `none`. | list |

### View a list of your lists

```
enbas show --type list
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `list`. | |

### View a specific list

Prints the information of the specified list to screen along with all the accounts added to it (if any).

```
enbas show --type list --list-id 01J1T9DWR20DC36QWZFKHWZJ3H
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `list`. | |
| `list-id` | string | false | The ID of the list you want to view. If this is not specified then a list of your lists will be printed instead. | |

### Edit a list

Edits the title and/or the replies policy of a list.

```
enbas edit --type list --list-id 01J1T9DWR20DC36QWZFKHWZJ3H --list-title "My Favourite People (in the world)"
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to edit.<br>Here this should be `list`. | |
| `list-title` | string | false | The title of the list that you want to edit. | |
| `list-replies-policy` | string | false | The policy of the replies for this list that you want to change to.<br>Valid values are `followed`, `list` and `none`. | |

### Delete a list

Deletes a list.

```
enbas delete --type list --list-id 01J1T9DWR20DC36QWZFKHWZJ3H
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to delete.<br>Here this should be `list`. | |
| `list-id` | string | true | The ID of the list you want to delete. | |

### Add accounts to a list

Adds one or more accounts to a list.

```
enbas add --type account --account-name @name@example.social --account-name @person@mastodon.example --to list --list-id 01J1T9DWR20DC36QWZFKHWZJ3H
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `account`. | |
| `account-name` | string | true | The name of the account you want to add to the list.<br>Use multiple times to specify multiple accounts. | |
| `to` | string | true | The resource you want to add the accounts to.<br>Here this should be `list`. | |
| `list-id` | string | true | The ID of the list that you want to add the accounts to. | |

### Remove accounts from a list

Removes one or more accounts from a list.

```
enbas remove --type account --account-name @person@mastodon.example --from list --list-id 01J1T9DWR20DC36QWZFKHWZJ3H
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `account`. | |
| `account-name` | string | true | The name of the account you want to remove from the list.<br>Use multiple times to specify multiple accounts. | |
| `to` | string | true | The resource you want to remove the accounts from.<br>Here this should be `list`. | |
| `list-id` | string | true | The ID of the list that you want to remove the accounts from. | |

## Timelines

### View a timeline

Prints a list of statuses from a timeline.

- View your home timeline.
   ```
   enbas show --type timeline --timeline-category home
   ```
- View a maximum of 5 statuses from your instance's public timeline.
   ```
   enbas show --type timeline --timeline-category public --limit 5
   ```
- View a timeline from one of your lists.
   ```
   enbas show --type timeline --timeline-category list --list-id 01J1T9DWR20DC36QWZFKHWZJ3H
   ```
- View a timeline from a hashtag.
   ```
   enbas show --type timeline --timeline-category tag --tag caturday
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `timeline`. | |
| `timeline-category` | string | false | The type of timeline you want to view.<br>Valid values are `home`, `public`, `list` and `tag`. | home |
| `list-id` | string | false | The ID of the list you want to view.<br>This is only required if `timeline-category` is set to `list`. | |
| `tag` | string | false | The hashtag you want to view.<br>This is only required if `timeline-category` is set to `tag`. | |
| `limit` | integer | false | The maximum number of statuses to print. | 20 |

## Media

### View media from a status

Downloads and opens media attachment(s) from a status.
Enbas currently supports viewing images and videos.
The media is downloaded to your cache directory before Enbas opens it with your preferred media player.
In order to view images and videos, you must specify your image viewer and
video player in your configuration file respectively.

See the [configuration reference page](configuration.md#integration) on how to set up integration with
your media players.

- View a specific media attachment from a specific status
   ```
   enbas show --type media --from status --status-id 01J0N11V4V7PWH0DDRAVT7TCFK --attachment-id 01J0N0RQSJ7CFGKHA30F7GBQXT
   ```
- View all image attachments from a specific status
   ```
   enbas show --type media --from status --status-id 01J0N11V4V7PWH0DDRAVT7TCFK --all-images
   ```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `media`. | |
| `from` | string | true | The resource you want to view the media from.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status that you want to view the media from. | |
| `attachment-id` | string | false | The ID of the media attachment to download and view.<br>Use this flag multiple times to specify multiple media attachments. | |
| `all-images` | boolean | false | Set to `true` to show all images from the status. | false |
| `all-videos` | boolean | false | Set to `true` to show all videos from the status. | false |

## Bookmarks

### View your bookmarks

Prints a list of statuses that you have bookmarked.

```
enbas show --type bookmarks
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to view.<br>Here this should be `bookmarks`. | |
| `limit` | integer | false | The maximum number of bookmarks to show. | 20 |

### Add a status to your bookmarks

```
enbas add --type status --status-id 01J17FH1KD9CN6J9Q01011NE0D --to bookmarks
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to add.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status. | |
| `to` | string | true | The resource you want to add the status to.<br>Here this should be `bookmarks`. | |

### Remove a status from your bookmarks

```
enbas remove --type status --status-id 01J17FH1KD9CN6J9Q01011NE0D --from bookmarks
```

| flag | type | required | description | default |
|------|------|----------|-------------|---------|
| `type` | string | true | The resource you want to remove.<br>Here this should be `status`. | |
| `status-id` | string | true | The ID of the status. | |
| `from` | string | true | The resource you want to remove the status to.<br>Here this should be `bookmarks`. | |

## Notifications

_Not yet supported_
