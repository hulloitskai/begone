# begone

_A fully automatic spamming tool, created for the sole purpose of
obliterating conversation threads on
[Facebook Messenger](https://messenger.com)._

[![grp-img]][grp]

It works with individual conversations as well as group threads—a real versatile
beast. Uses a modified version of
[`unixpickle/fbmsgr`](https://github.com/unixpickle/fbmsgr) as the underlying
Messenger client. And, obviously, written in [Go](https://golang.org).

<br />
<p align="center">
  <img src="./.github/demo.gif" width=600>
</p>
<br />

## Usage

### Installing

Grab the [latest release](https://github.com/stevenxie/begone/releases) compiled
for your system.

Ensure that the binary is executable, and place it somewhere in your `$PATH`.
For macOS users, this might look something like this:

```bash
$ mv ~/Downloads/begone-darwin-x64 /usr/local/bin/begone
$ chmod u+x /usr/local/bin/begone
```

### Running

```bash
## Save login credentials to ~/.begone.json.
$ begone login

## Launch an emoji attack on a conversation thread.
$ begone emojify https://messenger.com/t/exampleid
```

<br />

## Advanced Usage

### Making from source

> This requires the [Go](https://golang.org) language and associated toolchain
> to be installed. If you're on _macOS_, this may be as easy as `brew install go`!.

```bash
## Clone this repository.
$ git clone git@github.com:stevenxie/begone.git

## Compile and install a version for your machine.
$ make install  # (or go install)
```

## TODOs

- [ ] (maybe) Implement attacks using local files (images)?

[grp]: https://goreportcard.com/report/github.com/stevenxie/begone
[grp-img]: https://goreportcard.com/badge/github.com/stevenxie/begone
