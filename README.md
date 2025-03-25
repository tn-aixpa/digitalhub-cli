# CLI

The Command-Line Interface (CLI) is a simple tool that offers a number of functionalities to set the platform up.

The executable can be downloaded in the *Releases* page.

## Build from source

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

### `create`
`create` will create an instance of the indicated resource on the platform. It takes the following parameters:

- `-n environment` (Optional)
- `-p project` (Optional)
- `-e entity_type` (Optional)
- `yaml_file_path`

To create a project, the only mandatory parameter is the path to the yaml file containing its definition. To create other resources (artifacts, models, etc.), you will have to specify its type and the project it will belong to.

Create a project:
``` sh
./dhcli create samples/project.yaml
```

Create an artifact:
``` sh
./dhcli create -p my-project -e artifacts samples/artifact.yaml
```
#### Resource types
The type of a resource can be any of the following: ```artifacts```, ```dataitems```, ```models```, ```functions```, ```workflows```, ```runs```, ```secrets```.

### `read`
`read` returns a list of resources or the details of a specific resource. It takes the following parameters:

- `-n environment` (Optional)
- `-p project` (Optional)
- `-e entity_type` (Optional)
- `-i id` (Optional)

List all projects:
``` sh
./dhcli read
```

Read a project's details:
``` sh
./dhcli read -p my-project
```

List all artifacts in a project:
``` sh
./dhcli read -p my-project -e artifacts
```

Read an artifact's details:
``` sh
./dhcli read -p my-project -e artifacts -i 88b5cd516e334c0bbea7352a2aeb3fb9
```

### `update`
`update` will update a resource with a new definition. It takes the following parameters:

- `-n environment` (Optional)
- `-e entity_type` (Optional)
- `-i id` (Optional)
- `project`
- `yaml_file_path`

Update a project:
``` sh
./dhcli update my-project samples/project.yaml
```

Update an artifact:
``` sh
./dhcli update -e artifacts -i f884eb3fccca4ab0b701daaaf96358da my-project samples/artifact.yaml
```

### `delete`
`delete` will delete a resource. It takes the following parameters:

- `-n environment` (Optional)
- `-e entity_type` (Optional)
- `-i id` (Optional)
- `project`

Delete a project:
``` sh
./dhcli delete my-project
```

Delete an artifact:
``` sh
./dhcli delete -e artifacts -i my-artifact my-project
```