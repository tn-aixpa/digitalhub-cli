# dhcli

## `register`
`register` takes the following parameters:
- `name of section`
- `authorization provider`
``` sh
./dhcli register bologna aac.digitalhub-dev.smartcommunitylab.it
```
It will create a `.cli.ini` file in the user's home directory (or, if not possible, in the current one), generating a section with the configuration retrieved from the authorization provider.


## `login`
`login` is to be used after registering a section with the `register` command. It takes the following parameters:
- `name of section`
``` sh
./dhcli login bologna
```
It will read the corresponding section from the configuration file and log in to the authorization provider. It will update the section with the JWT token obtained.

## `logout`
`logout` takes the following parameters:
- `name of section`
``` sh
./dhcli logout bologna
```
It removes the JWT from the corresponding section in the configuration file.
