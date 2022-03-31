# SDM-EXT

The SDM-EXT project aims to create a CLI (Command Line Interface) so that the user can execute all the commands provided by the SDM CLI plus some exclusive commands from the SDM-EXT CLI. The "ext" stands for "extension".

# ****Installation and configuration****

You need to install SDM CLI and configure it on your machine. See [Installing the SDM CLI](https://www.strongdm.com/docs/user-guide/client-installation) for how to install. After that, make sure that the SDM CLI can be accessed using the sdm alias in your terminal.

To run the SDM-EXT CLI, after downloading the project and going to its directory, run the command `go build sdm-ext.go` and export the generated executable file to be able to run the SDM-EXT CLI using only the `sdm-ext` command.

# SDM-EXT flags

Using `sdm-ext` is basically the same as using the `sdm` command, but with some additional flags.

You can use the following sdm-ext flags:

- `--file value`: used to pass a JSON file;
- `-f value`: is also used to pass a JSON file;
- `--stdin`: used to get user input in JSON format;
- `-i`: is also used to get user input in JSON format;

# Examples

## Exclusive SDM-EXT CLI Commands

- `sdm-ext admin servers add --file ./json-examples/server-list.json`
  - It registers a list of servers in the strongDM app from a JSON file. 
- `sdm-ext admin servers add --f ./json-examples/server-list.json`
  - It also registers a list of servers in the strongDM app from a JSON file.
- `sdm-ext admin servers add --stdin`
  - It registers a list of servers in the strongDM app from a user input in JSON format.
- `sdm-ext admin servers add -i`
  - It also registers a list of servers in the strongDM app from a user input in JSON format.

## SDM CLI commands executed with the SDM-EXT CLI

- `sdm-ext admin servers list`
  - This uses the original SDM CLI to list all servers registered in the strongDM app.
