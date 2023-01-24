# boop

[![pipeline status](https://gitlab.com/genieindex/boop-go/badges/master/pipeline.svg)](https://gitlab.com/genieindex/boop-go/-/commits/master)

[boop.pdf](https://github.com/abedef/boop-go/files/10492741/boop.pdf)

boop is an experimental information capture and recall system, supporting syntax lightly inspired by [sowhat](https://github.com/tatatap-com/sowhat)

## Executables

## boop

Dependencies:
* an instance of [`boop-server`](https://gitlab.com/genieindex/boop-server)

### Installation

Install the latest version via Go:

```sh
go install gitlab.com/genieindex/boop-cli@latest
```

To refer to the command simply as `boop`, add `alias boop='boop-cli'` to your shell's initialization file.

### Usage

Assuming `boop-server` is deployed to [boop.example.com]() and configured for use with phone number `+15555551234`:

Create a configuration file (`~/.config/boop/config.yaml`, or `~/Library/Application Support/boop/config.yaml` on macOS) with the following contents:

```yaml
endpoint: http://boop.example.com/
phone: +15555551234
```

Now you can run `boop` to print all boops to stdout, or pipe text into `boop` like `echo hello world | boop` to save a boop.

## boopd

Dependencies:
* TODO

### Installation

TODO

### Usage

TODO
