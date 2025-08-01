# msc (Monktype's Stream Commands)

`msc` is a hacky command-line utility designed to perform various Twitch-related channel and moderation functions using the Twitch API.
Currently, it supports creating polls and sending announcements, with plans for further expansion.

## Installation

To install the application, clone the repository and build it using Go:

`git clone https://github.com/monktype/msc.git`
`cd msc`
`go build .`

(outputs binary in the current directory)

I may also try to keep binaries build in the "Releases" tab to the right -->

## Usage

After building the application, you can run it from the command line using the commands listed below or with `msc --help`.


## Setup

To use `msc`, you need to set up an application Twitch API key at `https://dev.twitch.tv` of type "public" (you can make the catagory "Application Integration" if you want).
Only the client ID is required for setup.

After obtaining your client ID, run the following command to set it up with the application:

`msc setup -i <client-id>`

By default, the application will attempt to authenticate after setup.
It will provide you with a link to authenticate against Twitch.
Your client-id is stored in your OS's keyring.
Once the OAuth authentication is completed, the OAuth token is also stored in your OS's keyring. 

Please note that the authentication key needs to be re-authenticated occasionally using:

`msc authenticate`

The exact duration for how long these keys last is still being determined, but it's less than a few days.
You can re-authenticate to get a fresh key before your previous key expires.

## Commands

### Version Command
Displays the current version of the application *if built with a tag*.

`msc version`

### Setup Command
Sets up the application with the required client ID from the Twitch Developer portal.

#### Flags:
- `-n`, `--no-auth`: Skip trying to authenticate after running setup.
- `-i`, `--client-id`: **(Required)** Client ID from the Twitch Dev portal.

#### Example:
See Setup section above.

### User ID Command
Retrieves the user ID associated with the account in the arguments.

#### Example:
`msc userid djclancy`

The above command returns `Username djclancy = ID 268669435`

### Poll Command
Creates a new poll in a specified channel. The command requires standalone string arguments (between 2 to 10).

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

## Contributing

Feel free to submit issues or pull requests to improve the project.
I'm still working on expanding the functions in this!

## License

This project is licensed under the MIT License.

