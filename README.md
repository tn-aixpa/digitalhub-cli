# dhcli

## `register`
`register` takes the following positional parameters:
- `environment_name`
- `authorization_provider`
- `client_id`
``` sh
./dhcli register bologna aac.digitalhub-dev.smartcommunitylab.it c_dhcliclientid
```
It will create a `.cli.ini` file in the user's home directory (or, if not possible, in the current one), generating a section with the specified environment name and containing the configuration retrieved from the authorization provider.

## `use`
`use` takes the following positional parameters:
- `default environment_name`
``` sh
./dhcli use bologna
```
It sets the environment to use when none is specified in the configuration file's default section.

## `login`
`login` is to be used after registering an environment with the `register` command. It takes the following optional parameter:
- `-e environment_name`
``` sh
./dhcli login -e bologna
```
It will read the corresponding section from the configuration file and log in to the authorization provider. It will update the section with the access token obtained.

## `remove`
`remove` takes the following positional parameters:
- `environment_name`
``` sh
./dhcli remove bologna
```
It removes the section from the configuration file.
