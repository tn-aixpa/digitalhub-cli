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

Available commands and their parameters are listed here. In these examples, the executable is named `dhcli`. When you provide flag parameters, make sure they are listed **before** positional ones.

### `register`
`register` takes the following parameters:

- `-e environment` *Optional*. Name of the environment to register.
- `core_endpoint` **Mandatory**

``` sh
dhcli register -e example http://localhost:8080
```
It will create a `.dhcore.ini` file (if it doesn't already exist) in the user's home directory, or, if not possible, in the current one. A section will be appended, using the provided environment name (or, if missing, the one returned by the endpoint), containing the environment's configuration. This environment will be set as default, unless one is already set.

### `list-env`
`list-env` lists available environments. It takes no parameters.

``` sh
dhcli list-env
```

### `use`
`use` takes the following parameters:

- `environment` **Mandatory**

``` sh
dhcli use example
```
This will set the default environment.

### `login`
`login` is to be used after registering an environment with the `register` command. It takes the following parameters:

- `environment` *Optional*.

``` sh
dhcli login example
```
It will read the corresponding section from the configuration file and start the log in procedure. It will update this section with the access token obtained. If no environment is specified, it will use the default one.

### `refresh`
`refresh` is to be used after the `login` command, to update `access_token` and `refresh_token`. It takes the following parameters:

- `environment` *Optional*

``` sh
dhcli refresh example
```
If no environment is specified, it will use the default one.

### `remove`
`remove` takes the following parameters:

- `environment` **Mandatory**

``` sh
dhcli remove example
```
It will remove the section from the configuration file.

### `init`
`init` is used to install the platform's Python packages; therefore, Python must be installed. It takes the following parameters:

- `environment` *Optional*

``` sh
dhcli init example
```
It will match core's minor version as indicated in the specified environment. If no environment is specified, it will use the default one.

### `create`
`create` will create an instance of the indicated resource on the platform. It takes the following parameters:

- `-e environment` *Optional*
- `-p project` *Optional* (ignored) when creating projects, **mandatory** otherwise.
- `-f yaml_file_path` **Mandatory** when creating resources other than projects, *alternative* to `name` for projects.
- `-n name` *Optional* (ignored) when creating resources other than projects, *alternative* to `yaml_file_path` for projects.
- `-reset-id` *Optional*. Boolean. If set, the `id` specified in the file is ignored.
- `resource` **Mandatory**

The type of resource to create is mandatory. The project flag `-p` is only mandatory when creating resources other than projects (artifacts, models, etc.). For projects, you may omit the file path and just use the `-n` flag to specify the name. The `-reset-id` flag, when set, ensures the created object has a randomly-generated ID, ignoring the `id` field if present in the input file (this is not relevant to projects).

Create a project:
``` sh
dhcli create -f samples/project.yaml projects
```

Create an artifact, while resetting its ID:
``` sh
dhcli create -p my-project -f samples/artifact.yaml -reset-id artifacts
```

#### Resource types
The `resource` positional parameter can accept any value (to support future updates), but if an invalid one is specified, the CLI will forward the error returned by core. This parameter is used in building the endpoint of the URL for the API call to core's back-end, therefore, it is possible to specify aliases for a resource in the `config.json` file.

### `list`
`list` returns a list of resources of the specified type. It takes the following parameters:
- `-e environment` *Optional*
- `-o output_format` *Optional*. Accepts `short`, `json`, `yaml`. Defaults to `short`.
- `-p project` *Optional* (ignored) for projects, **mandatory** otherwise.
- `-n name` *Optional*. If present, will return all versions of specified resource. If missing, will return the latest version of all matching resources.
- `-k kind` *Optional*
- `-s state` *Optional*
- `resource` **Mandatory**

`output_format` determines how the output will be formatted. The default value, `short`, is meant to be used to quickly check resources in the terminal, while `json` and `yaml` will format the output accordingly, making it ideal to write to file.

List all projects:

``` sh
dhcli list projects
```

List all artifacts in a project:

``` sh
dhcli list -p my-project artifacts
```

Note that you can easily write the results to file by redirecting standard output:
``` sh
dhcli list -o yaml -p my-project artifacts > output.yaml
```

### `get`
`get` returns the details of a single resource. It takes the following parameters:
- `-e environment` *Optional*
- `-o output_format` *Optional*. Accepts `short`, `json`, `yaml`. Defaults to `short`.
- `-p project` *Optional* (ignored) for projects, **mandatory** otherwise.
- `-n name` *Optional* (ignored) if `id` is missing, **mandatory** otherwise. Will return latest version of specified resource.
- `resource` **Mandatory**
- `id` *Optional* if `-n name` is missing, **mandatory** otherwise.

Similarly to the `list` command, `output_format` determines how the output will be formatted. The default value, `short`, is meant to be used to quickly check resources in the terminal, while `json` and `yaml` will format the output accordingly, making it ideal to write to file.

Get project:

``` sh
dhcli get projects my-project
```

Get artifact:
``` sh
dhcli get -p my-project artifacts my-artifact-id
```

Get artifact and write to file:
``` sh
dhcli get -o yaml -p my-project artifacts my-artifact-id > output.yaml
```

### `update`
`update` will update a resource with a new definition. It takes the following parameters:

- `-e environment` *Optional*
- `-p project` *Optional* (ignored) for projects, **mandatory** otherwise.
- `-f yaml_file_path` **Mandatory**
- `resource` **Mandatory**
- `id` **Mandatory**

Update a project:
``` sh
dhcli update -f samples/project.yaml projects my-project
```

Update an artifact:
``` sh
dhcli update -p my-project -f samples/artifact.yaml artifacts my-artifact-id
```

### `delete`
`delete` will delete a resource. It takes the following parameters:

- `-e environment` *Optional*
- `-p project` *Optional* (ignored) for projects, **mandatory** otherwise.
- `-n name` *Alternative* to `id`, will delete all versions of a resource.
- `-y` *Optional*. Boolean. If omitted, confirmation will be asked.
- `-c` *Optional*. Boolean, only applies to projects. When set, all resource belonging to the project will also be deleted.
- `resource` **Mandatory**
- `id` *Alternative* to `name`, will delete a specific version. For projects, since versions do not apply, this is synonym with `id`.

Delete a project and all of its resources:
``` sh
dhcli delete -c projects my-project
```

Delete an artifact, skip confirmation:
``` sh
dhcli delete -p my-project -y artifacts my-artifact-id
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
