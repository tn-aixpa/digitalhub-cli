# CLI

The Command-Line Interface (CLI) is a simple tool that offers a number of functionalities to set the platform up.

The executable can be downloaded in the *Releases* page.

## Installation

Download and unpack the corresponding archive from the [GitHub release page](https://github.com/scc-digitalhub/digitalhub-cli/releases). Run the ``dhcli`` executable. 

In case of Mac, it is possible to use Homebrew Tap distribution:

``` sh
brew tap scc-digitalhub/digitalhub-cli https://github.com/scc-digitalhub/digitalhub-cli
brew install dhcli
```

## Commands

Available commands and their parameters are listed here. In these examples, the executable is named `dhcli`. When you provide *optional* parameters, make sure they are listed **before** *mandatory* ones.

### `register`
`register` takes the following parameters:

- `-n name` (Optional)
- `core_endpoint`

``` sh
./dhcli register -n example http://localhost:8080
```
It will create a `.dhcore.ini` file (if it doesn't already exist) in the user's home directory, or, if not possible, in the current one. A section will be appended, using the provided name (or, if missing, the one returned by the endpoint), containing the environment's configuration. This environment will be set as default, unless one is already set.

### `use`
`use` takes the following parameters:

- `environment`

``` sh
./dhcli use example
```
This will set the default environment.

### `login`
`login` is to be used after registering an environment with the `register` command. It takes the following parameters:

- `environment` (Optional)

``` sh
./dhcli login example
```
It will read the corresponding section from the configuration file and start the log in procedure. It will update this section with the access token obtained. If no environment is specified, it will use the default one.

### `refresh`
`refresh` is to be used after the `login` command, to update `access_token` and `refresh_token`. It takes the following parameters:

- `environment` (Optional)

``` sh
./dhcli refresh example
```
If no environment is specified, it will use the default one.

### `remove`
`remove` takes the following parameters:

- `environment`

``` sh
./dhcli remove example
```
It will remove the section from the configuration file.

### `init`
`init` takes the following parameters:

- `environment` (Optional)

``` sh
./dhcli init example
```
It will install the python package through pip, matching core's minor version as indicated in the specified environment. If no environment is specified, it will use the default one.


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