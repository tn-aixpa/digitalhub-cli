# CLI

The Command-Line Interface (CLI) is a simple tool that offers a number of functionalities to set the platform up.

The executable can be downloaded in the *Releases* page.

For a full overview of available commands, please check [the documentation](https://scc-digitalhub.github.io/docs/0.11/cli/commands/).

## Installation

Download and unpack the corresponding archive from the [GitHub release page](https://github.com/scc-digitalhub/digitalhub-cli/releases). Run the ``dhcli`` executable. 

On Mac, you can use Homebrew Tap distribution:

``` sh
brew tap scc-digitalhub/digitalhub-cli https://github.com/scc-digitalhub/digitalhub-cli
brew install dhcli
```

## Build and publish

### Build and publish with GoReleaser

Install the [GoReleaser](https://goreleaser.com/install/) tool.

Create a corresponding tag for the code

``` sh
git tag -a X.0.0 -m "Some release" && git push origin X.0.0
```

Make a relase that will 

- generate the artifacts for Linux (i386, amd64, arm64), Windows (i386, amd64, arm64) and Mac (universal package)
- be published to GitHub and to Homebrew Tap 

``` sh
export GITHUB_TOKEN=YOURGITHUBTOKENHERE; goreleaser --clean
```

### Build from source

If you wish to build the executable from source, run the following:

``` sh
go build
```

It will generate an executable named `dhcli` for your operating system and architecture. To change the target OS or architecture, you need to set the `GOOS` and `GOARCH` variables and build it. Some examples:
``` sh
GOOS=linux GOARCH=amd64 go build -o dhcli-linux-amd64
```
``` sh
GOOS=darwin GOARCH=arm64 go build -o dhcli-darwin-arm64
```
``` sh
GOOS=windows GOARCH=amd64 go build -o dhcli-win-amd64.exe
```

For a complete list of available values, run `go tool dist list`, which will return the list of valid OS/ARCH combinations.
