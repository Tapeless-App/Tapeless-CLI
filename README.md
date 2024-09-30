# Welcome to the `Tapeless` CLI

The `Tapeless-CLI` is a complementary tool for the [Tapeless web application](https://tapeless.app).

# User Guide

## Quickstart

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

First, setup the version you wish to build. e.g. `TAPELESS_VERSION=1.0.1` and run the build command:

```
  go build -o ./build/tapeless -ldflags "\
    -X tapeless.app/tapeless-cli/env.Version=${TAPELESS_VERSION} \
    -X tapeless.app/tapeless-cli/env.ApiURL=https://api.tapeless.app/cli \
    -X tapeless.app/tapeless-cli/env.WebURL=https://tapeless.app \
    -X tapeless.app/tapeless-cli/env.LoginCallbackPort=8080"
```
