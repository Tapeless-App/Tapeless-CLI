# Welcome to the `Tapeless` CLI

The `Tapeless-CLI` is a complementary tool for the [Tapeless web application](https://tapeless.app).

# User Guide

## Quick Start

### Installation

**MacOs & Linux**

The Tapeless CLI is best installed via Homebrew:

```
brew install Tapeless-App/Tapeless-CLI/tapeless
```

This will install the pre-built binaries. To build from source, follow the "developer guide" below.

**Windows**

Download the correct windows artifact from the [latest releases](https://github.com/Tapeless-App/Tapeless-CLI/releases/latest), extract the `.exe` and add it to your `PATH`.

### Usage

1. `tapeless status` check your setup and see what action to perform next
1. `tapeless version` verify installation and ensure you are using the latest version
1. `tapeless login` will setup your session by logging into tapeless via the web UI
1. `tapeless projects add` allows you to create a new projects
1. `tapeless repos add` will add the working directory as a git repository to one of your projects
1. `tapeless sync` will push the git commits of all your registered repositories to the respective projects on tapeless
1. `tapeless open` will open the [Tapeless web application](https://tapeless.app) in your default browser

## Local Config

By default, all of your local config is stored in `[HOME_DIR]/.tapeless/config`. It is highly recommended **NOT** to edit the configuration manually, but instead use the Tapeless CLI commands.

# Developer Guide

## Running Locally

### Setup

Ensure that the variables in [./env/env.go](./env/env.go) are properly configured. You can run the CLI either with a locally running Tapeless instance, or with the production build, assuming the necessary endpoints are deployed.

### Running

Simply run `go run main.go [COMMAND]` to execute the desired command.

## Building

### Local Build

To create a local build, setup the version you wish to build. e.g. `TAPELESS_VERSION=1.0.0` and run the build command:

```
  go build -o ./build/tapeless -ldflags "\
    -X tapeless.app/tapeless-cli/env.Version=${TAPELESS_VERSION} \
    -X tapeless.app/tapeless-cli/env.ApiURL=https://api.tapeless.app/cli \
    -X tapeless.app/tapeless-cli/env.WebURL=https://tapeless.app \
    -X tapeless.app/tapeless-cli/env.LoginCallbackPort=8080"
```

### Release Build

This project uses GoReleaser to create binaries for all systems and associate them with the latest release tag.

First make sure all changes are committed and create a new release, e.g.

```
git tag -a v0.0.4 -m "Description for the release..."
```

Then, create your binaries and push the release via:

```
GITHUB_TOKEN=[YOUR_GITHUB_TOKEN] goreleaser release
```

Finally, update the formula in the [homebrew-tapeless-cli](https://github.com/Tapeless-App/homebrew-tapeless-cli) repository with the latest version and shas.
