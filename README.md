# stackup-paymaster

A JSON-RPC client for serving [ERC-4337 verifying paymaster requests](https://docs.stackup.sh/docs/paymaster-api-rpc-methods).

> **⚠️ This software is still in early development. It is not yet recommended for business critical systems.**

# Running an instance

```
Documentation coming soon...
```

# Contributing

Steps for setting up a local dev environment for contributing to the Paymaster.

## Prerequisites

- Go 1.21 or later
- [ERC-4337 Devnet](https://github.com/stackup-wallet/erc-4337-devnet) running

## Setup

```bash
# Installs https://github.com/cosmtrek/air for live reloading.
# Runs go mod tidy.
# Creates a .env file with sensible default values for local development.
make install-dev
```

## Run the paymaster service

Start a local paymaster instance:

```bash
make dev-run
```

# License

Distributed under the GPL-3.0 License. See [LICENSE](./LICENSE) for more information.

# Contact

Feel free to direct any technical related questions to the `dev-hub` channel in the [Stackup Discord](https://discord.gg/VTjJGvMNyW).
