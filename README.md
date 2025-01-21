# dhcli

Before interacting with the CLI, you must register a corresponding app in AAC. You can use the contents of `sample-client.yaml` to import an app quickly.

## `register`
`register` takes the following parameters:
- `-s scope` (Optional)
- `environment_name`
- `authorization_provider`
- `client_id`
``` sh
./dhcli register -s offline_access bologna aac.digitalhub-dev.smartcommunitylab.it c_dhcliclientid
```
It will create a `.cli.ini` file in the user's home directory (or, if not possible, in the current one), generating a section with the specified environment name and containing the configuration retrieved from the authorization provider.

## `use`
`use` takes the following parameters:
- `default_environment_name`
``` sh
./dhcli use bologna
```
It sets the environment to use when none is specified in the configuration file's default section.

## `login`
`login` is to be used after registering an environment with the `register` command. It takes the following parameters:
- `-e environment_name` (Optional)
``` sh
./dhcli login -e bologna
```
It will read the corresponding section from the configuration file and log in to the authorization provider. It will update the section with the access token obtained. If no environment is specified, it will use the one set by the `use` command.

## `refresh`
`refresh` is to be used after the `login` command, to update `access_token` and `refresh_token`. It takes the following parameters:
- `-e environment_name` (Optional)
``` sh
./dhcli refresh -e bologna
```
If no environment is specified, it will use the one set by the `use` command.

## `remove`
`remove` takes the following parameters:
- `environment_name`
``` sh
./dhcli remove bologna
```
It removes the section from the configuration file.
