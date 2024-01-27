# Description

A project whose purpose is to build bitwarden-cli for various platforms and package it in various ways. 

Supported platforms are:

- darwin/arm64
- darwin/amd64
- linux/arm64
- linux/amd64

Supported packaging is:

- standalone binary
- docker image

# Getting started

1. create a `.env` to store any sensitive environment var and export them using: `export $(xargs <.env)`
2. `dagger run go run main.go` to get CLI usage 

# Context

- bitwarden-cli teams doesn't provide such platforms/packages
- opportunity to practice go, cobra/viper, dagger.io and github ecosystem

# Design

- go program running dagger client
- using cobra/viper for CLI UI

# Roadmap

- detector when new release upstream published
- ability to run manual pipeline and set env var to force release version
- golang ginkgo/unit test
- use act locally
- config github action avec dagger
- tester github action depuis vscode
- config validargs constraints avec viper (ie: pour nodejs value limited choice)
- github action when upstream release new version (with dependabot)
- use github API to retrieve zipfile
- tool pre-commit pour verif si fuite info sensible
- cache npm/deb,... content
- make registry auth optional

# HowTo

## Add command

```
cobra-cli add lint --viper
```