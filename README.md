# Snek
Snek is a CLI tool for interacting with [Longboi](https://github.com/che-ict/longboi) servers, built with [Golang](https://golang.org/), [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper).

## Installation
Binary builds are provided for Windows, Mac and Linux. See the [releases page](https://github.com/che-ict/snek/releases) for more information.
A small install script is also provided for Windows:
```ps
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/che-ict/snek/main/install.ps1'))
```

## Usage
Snek has a few commands:

### Auth
You can use the `auth login` command to authenticate with a Longboi server. Credentials are stored in a configuration file. You can validate your credentials by using the `auth validate` command.

Using `auth login {server} --web`, you can log in using the web interface, so you don't have to manually create an API key

### Checkout
You can use the `checkout` command to checkout a project from a Longboi server. You can find available courses on the server by visiting its web interface.

### Pull
In a project directory, you can use the `pull` command to pull the latest changes from the server.

### Attempt
In an exercise directory, you can use the `attempt` command to attempt a submission.

### Version
Prints the current version of the program

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.