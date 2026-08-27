# msc (Monktype's Stream Commands)

`msc` is a hacky command-line utility designed to perform various Twitch-related channel and moderation functions using the Twitch API.
Currently, it supports creating polls, sending announcements, and more (see below) with plans for further expansion.

## Installation

To install the application, clone the repository and build it using Go:

```
git clone https://github.com/monktype/msc.git
cd msc
go build .
```

(outputs binary in the current directory)

I may also try to keep binaries built in the "Releases" tab to the right -->

## Usage

After building the application, you can run it from the command line using the commands listed below or with `msc --help`.


## Setup

To use `msc`, you need to set up an application Twitch API key at `https://dev.twitch.tv`.

There are two different methods depending on your setup:

#### Authorization Code Grant Flow
Pros: Re-authentication happens automatically / requires less frequent human re-authentication.

Cons: Requires copying in a secret, too.

Create a key of type "confidential" (you can make the catagory "Application Integration" if you want).
Redirect to `http://localhost:3024/redirect`.

Both the client ID and a client secret are necessary for this configuration.

After obtaining the client ID and the client secret, run the following command to set it up with the application:

`msc setup -i <client-id> -s`

You will be prompted within the application to input the secret.

By default, the application will attempt to authenticate after setup.
It will provide you with a link to authenticate against Twitch.
Your client-id, secret, refresh token, etc are all stored in your OS's keyring.

#### Implicit Grant Flow
Pros: Only a client ID needs to be copied into the client.

Cons: Re-authentication needs to happen about once a day.

Create a key of type "public" (you can make the catagory "Application Integration" if you want).
Redirect to `http://localhost:3024/redirect`.
Only the client ID is required for setup.

After obtaining your client ID, run the following command to set it up with the application:

`msc setup -i <client-id>`

By default, the application will attempt to authenticate after setup.
It will provide you with a link to authenticate against Twitch.
Your client-id is stored in your OS's keyring.
Once the OAuth authentication is completed, the OAuth token is also stored in your OS's keyring. 

Please note that the authentication key needs to be re-authenticated occasionally using:

`msc authenticate`

The duration of these access tokens from Twitch are about one day.
You can re-authenticate to get a fresh key before your previous key expires.
Errors do not necessarily get returned when the token has expired.

**Note (new scopes):** the set of OAuth scopes `msc` requests has been expanded
(moderation scopes for ban management and scopes for predictions were added).
You'll need to re-authenticate **once** — run `msc authenticate`, or re-run `msc setup` —
so the token stored in your keyring picks up all of the new scopes at once.
Subsequent token refreshes will keep those scopes; you do not need to re-auth for each one.

## Commands

### Version Command
Displays the current version of the application *if built with a tag*.

`msc version`

### Setup Command
Sets up the application with the required client ID and any other required components (see above) from the Twitch Developer portal.

#### Flags:
- `-n`, `--no-auth`: Skip trying to authenticate after running setup.
- `-i`, `--client-id`: **(Required)** Client ID from the Twitch Dev portal.
- `-s`, `--secret`: Add a secret for code authentication instead of token authentication.

#### Example:
See Setup section above.

### User ID Command
Retrieves the user ID associated with the account in the arguments.

#### Example:
`msc userid djclancy`

The above command returns `Username djclancy = ID 268669435`

### Poll Command
Creates a new poll in a specified channel. The command requires standalone string arguments (between 2 to 5).

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.
- `-d`, `--duration`: **(Required)** Duration in seconds.
- `-t`, `--title`: **(Required)** Title for the poll.
- `-n`, `--no-watch`: Skip watching until the end of the poll, just print the poll ID.
- `-a`, `--send-announcement`: Send an announcement when the poll starts.
- `-A`, `--send-announcement-result`: Send an announcement when the poll starts AND when the poll ends with the result. Implies `-a`.

#### Example:
`msc poll -c djclancy -d 15 -t "Yes or no?" "Yes" "No"`

This creates a 15-second poll on djclancy's channel with "Yes" and "No" as options.

### Prediction Command
Creates a new prediction in a specified channel. The command requires exactly 2 standalone string arguments (the two outcome choices). A prediction does not close itself — once created, end it with the `resolve` subcommand (pick a winner) or the `cancel` subcommand (refund channel points). Use the `lock` subcommand to close the prediction window early (stop accepting new predictions before the timer runs out, without picking a winner or refunding).

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.
- `-d`, `--window`: **(Required)** Prediction window in seconds (1..1800).
- `-t`, `--title`: **(Required)** Title for the prediction.

#### Subcommands:
- `resolve`: Pick the winning outcome (channel points are awarded to those who predicted it).
  - `-c`, `--channel-name`: **(Required)** Target channel name.
  - `-i`, `--prediction-id`: **(Required)** Prediction ID (printed by the create command).
  - `-o`, `--outcome`: **(Required)** Winning outcome: `1` or `2`, or the outcome's title.
- `cancel`: Cancel the prediction with no winner (channel points are refunded).
  - `-c`, `--channel-name`: **(Required)** Target channel name.
  - `-i`, `--prediction-id`: **(Required)** Prediction ID.
- `lock`: Lock a running prediction early, closing the window so no more predictions are accepted (no winner, no refund; it can still be resolved afterward).
  - `-c`, `--channel-name`: **(Required)** Target channel name.
  - `-i`, `--prediction-id`: **(Required)** Prediction ID.

#### Examples:
`msc prediction -c djclancy -d 30 -t "Will the next clip be a highlight?" "Yes" "No"`

This creates a 30-second prediction on djclancy's channel with "Yes" and "No" as the two outcomes. It prints the prediction ID, the numbered outcomes, and the exact `resolve`/`cancel`/`lock` commands to close it out.

`msc prediction resolve -c djclancy -i <prediction-id> --outcome Yes` (pick the "Yes" outcome as the winner)

`msc prediction resolve -c djclancy -i <prediction-id> --outcome 2` (pick the second outcome by index)

`msc prediction cancel -c djclancy -i <prediction-id>` (cancel, refunding channel points)

`msc prediction lock -c djclancy -i <prediction-id>` (lock the window early, stopping new entries before the timer ends)

### Announcement Command
Sends an announcement to a specified channel. Every string argument is passed as text in the announcement.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.
- `-b`, `--border-color`: Border color (default is "primary"). Options: primary, blue, green, orange, purple.

#### Example:
`msc announcement --channel-name djclancy --border-color blue "This is an announcement!"`

### Shoutout Command
Shoutout a specified channel on a specified channel.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name (location of shoutout).
- `-s`, `--shoutout-name`: **(Required)** Channel to shoutout.

### Start-Ad Command
Start advertisements / commercials on a specific channel.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.
- `-l`, `--length`: Length in seconds: 30, 60, 90, 120, 150, or 180; defaults to 60.

### Emote Only Mode Commands
Two commands related to Emote Only mode:
- `on`: Turn on Emote Only mode.
- `off`: Turn off Emote Only mode.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.

#### Examples:
`msc emote-only -c djclancy on`

`msc emote-only -c djclancy off`

### Subscribers Only Mode Commands
Two commands related to Subscribers Only mode:
- `on`: Turn on Subscribers Only mode.
- `off`: Turn off Subscribers Only mode.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.

#### Examples:
`msc submode -c djclancy on`

`msc submode -c djclancy off`

### Follower Only Mode Commands
Three commands related to Followers Only mode:
- `on`: Turn on Follower Only mode.
- `off`: Turn off Follower Only mode.
- `duration`: Turn on Follower Only mode (if off) with a specified duration in minutes.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.
- `-d`, `--duration`: **(Required for `duration`)** Duration in minutes (0..129600 valid).

#### Examples:
`msc follower-only -c djclancy on`

`msc follower-only -c djclancy off`

`msc follower-only -c djclancy duration -d 15` (turns on Follower Only mode on djclancy's channel with a 15 minute wait time)

### Slow Mode Commands
Three commands related to Slowmode:
- `on`: Turn on Slowmode.
- `off`: Turn off Slowmode.
- `duration`: Turn on Slowmode (if off) with a specified duration in seconds.

#### Flags:
- `-c`, `--channel-name`: **(Required)** Target channel name.
- `-d`, `--duration`: **(Required for `duration`)** Duration in seconds (3..120 valid).

#### Examples:
`msc slowmode -c djclancy on`

`msc slowmode -c djclancy off`

`msc slowmode -c djclancy duration -d 15` (turns on Slowmode on djclancy's channel with a 15 second chat cooldown)

### Channel Points Custom Redeems Commands
Six commands related to Channel Poitns Custom Redeems:
- `cancel`: Cancel a redemption instance, refunding the user.
- `create`: Create a channel point reward.
- `delete`: Delete a channel point reward.
- `fulfill`: Fulfill a redemption instance.
- `get`: Get Channel Point Rewards for channel.
- `redemptions`: Get channel point redemptions.

#### Flags:
Many; see `--help` for each subcommand above.

#### Examples:
`msc reward create -c djclancy -t "Poetry Slam" -p 1000 -i -u "Write a poem here and I'll read it out loud."`

`msc reward get -c djclancy`

`msc reward redemptions -c djclancy -r 25b0b2e2-7800-407c-a52b-9864ba6f6565 -s UNFULFILLED`

`msc reward fulfill -c djclancy -r 25b0b2e2-7800-407c-a52b-9864ba6f6565 -i 1f8ae074-28b6-428e-b745-dc25903848c8`

`msc reward delete -c djclancy -r 25b0b2e2-7800-407c-a52b-9864ba6f6565`

### Ban Commands
Three commands related to banning and timing out users. Note that the "banned" list returned by the API also includes timed-out users (a timeout is a ban with an expiry, while a perma-ban has no expiry).
- `list`: List banned and timed-out users.
- `add`: Ban or time out a user.
- `remove`: Remove a ban or timeout.

#### Flags:
- `list`: `-c`, `--channel-name` **(Required)**; optional `-a`, `--after` pagination cursor; optional `-A`, `--all` to fetch every page and list all banned/timed-out users (ignores `-a`).
- `add`: `-c`, `--channel-name` **(Required)**; give either `-u`, `--user` (username) or `--user-id` (raw ID); `-d`, `--duration` timeout in seconds (0 or omitted = perma-ban); `-r`, `--reason` ban reason.
- `remove`: `-c`, `--channel-name` **(Required)**; give either `-u`, `--user` (username) or `--user-id` (raw ID, for when the username no longer resolves).

#### Examples:
`msc ban list -c djclancy`

`msc ban list -c djclancy -A` (list every banned/timed-out user across all pages)

`msc ban add -c djclancy -u offender -d 600 -r "Spamming"` (600-second timeout)

`msc ban add -c djclancy -u offender` (perma-ban, since no duration is given)

`msc ban add -c djclancy --user-id 268669435 -d 600` (time out by ID — use when the username no longer resolves)

`msc ban remove -c djclancy -u offender`

`msc ban remove -c djclancy --user-id 268669435` (remove a ban by ID — use when the user no longer resolves by name; grab the ID from `ban list`'s `[id]` field)

Note: there is no batch-unban in the Twitch API, so removing many users means one (rate-limited) call per user.

## Contributing

Feel free to submit issues or pull requests to improve the project!

### Probable Next Additions
- Blocked Terms

## License

This project is licensed under the MIT License.

