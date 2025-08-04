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

## Contributing

Feel free to submit issues or pull requests to improve the project!

If you plan to poke at this project, take note that this project is using a [temporary forked library](https://github.com/Monktype/helix) until [an upstream PR](https://github.com/nicklaw5/helix/pull/244) is merged.

### Probable Next Additions
- Blocked Terms
- Predictions

## License

This project is licensed under the MIT License.

