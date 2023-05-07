# Simple Control and Command Server

This is a small project for my friend as he sails around the world.

Basically extends a CLI to be partially available over SMS.

## Requirements

- Golang >=1.20
- ngrok `brew install ngrok/ngrok/ngrok`
- Twilio Account
- yq (yaml parser)

## Installation

Use `make install`. Follow the instructions, then execute `cnc help` for a list
of available commands

### ngrok

Follow the [getting started guide](https://ngrok.com/docs/getting-started/)

## Usage

This is structure as a CLI, run `cnc help` after installation to see options.

When launching this as a server for handling SMS messages, execute `make run-ingress`.

Add the output forwarding URL as a Webhook in the Twilio console.

## Development

The whole point of this was to make the command extensible. In order to add
new commands, just follow the cobra documentation 
and create new `<cmd>.go` files in `cmd/`.

If you wish certain commands to be disallowed from internal usage, check
for that flag being set as a PersistentPreRun set in the root of the subcommand.

This process will create a new CLI entry point (for the developer), and
a new external entry point (for the sailor).

## TODO

- Check Twilio timeout, see if this can be increased
- Allow for multiple messages to be a query, currently incoming sat phones are limited to 55 chars
- Add news target
- Add rate limiting to the SMS server
- Add chained target
- Add option to persist to SQLite file
- Use more idiomatic installation method
- Add [validators](https://github.com/go-validator/validator) for inputs
