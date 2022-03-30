# SDM-EXT

The SDM-EXT project aims to create a CLI (Command Line Interface) so that the user can execute all the commands provided by the SDM CLI plus some exclusive commands from the SDM-EXT CLI. The "ext" stands for "extension".

# ****Installation and configuration****

You need to install SDM CLI and configure it on your machine. See [Installing the SDM CLI](https://www.strongdm.com/docs/user-guide/client-installation) for how to install. After that, make sure that the SDM CLI can be accessed using the sdm alias in your terminal.

To use the SDM-EXT CLI, you need to provide the following mandatory variables:

- SDM_API_ACCESS_KEY: SDM API access key;
- SDM_API_SECRET_KEY: SDM API Access Key Secret.

Ps.: you can use .env.example to create your .env file with variables.

![Screenshot from 2022-03-30 10-34-23](https://user-images.githubusercontent.com/49795183/160855074-0dee81d6-fa8d-4292-9d6d-04b3d76f773f.png)

To run the SDM-EXT CLI, after downloading the project, going to its directory and setting the variables above, run the command `go build sdm-ext.go` and export the generated executable file to be able to run the SDM-EXT CLI using only the `sdm-ext` command.

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
- `sdm-ext admin servers add --f ./json-examples/server-list.json`
- `sdm-ext admin servers add --stdin`
- `sdm-ext admin servers add -i`

## SDM CLI commands executed with the SDM-EXT CLI

- `sdm-ext admin servers list`
