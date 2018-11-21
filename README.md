# begone

_A fully automatic spamming tool, created for the sole purpose of
obliterating conversation threads on
[Facebook Messenger](https://messenger.com)._

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

## Launch an emoji attack on this conversation thread.
$ begone emoji https://messenger.com/t/exampleid
```

<br />

## Advanced Usage

### Making from source

> This requires the [Go](https://golang.org) language and associated toolchain
> to be installed. If you're on _macOS_, this may be as easy as `brew install go`!.

```bash
## Clone this repository.
$ git clone git@github.com:stevenxie/begone

## Compile and install a version for your machine.
$ make install  # (or go install)
```
