# Kinu no Michi (Silk Road)

## Setup on Linux/macOS/Windows with WSL2

You're going to need to sign up for [CoinAPI](https://www.coinapi.io/) to get an API token for use with their exchange rates API.

```shell
# Install go and node using mise (https://mise.jdx.dev)
mise install

# Install sqlc and goose with Homebrew (macOS/Linux)
brew bundle

# Install node packages
pnpm install

# Install Go dependencies and tooling
go mod download
go install tool

# Use sample .envrc to set environment variables
cp .envrc.sample .envrc
direnv allow

# Create and migrate database
createdb
goose up
```
