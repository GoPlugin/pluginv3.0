# Plugin Smart Contracts

## Installation

```sh
# via pnpm
$ pnpm add @plugin/contracts
# via npm
$ npm install @plugin/contracts --save
```

### Directory Structure

```sh
@plugin/contracts
├── src # Solidity contracts
│   └── v0.8
└── abi # ABI json output
    └── v0.8
```

### Usage

The solidity smart contracts themselves can be imported via the `src` directory of `@plugin/contracts`:

```solidity
import '@plugin/contracts/src/v0.8/AutomationCompatibleInterface.sol';
```

## Local Development

Note: Contracts in `dev/` directories are under active development and are likely unaudited. Please refrain from using these in production applications.

```bash
# Clone Plugin repository
$ git clone https://github.com/goplugin/pluginv3.0.git
# Continuing via pnpm
$ cd contracts/
$ pnpm
$ pnpm test
```

## Contributing

Please try to adhere to [Solidity Style Guide](https://github.com/goplugin/pluginv3.0/blob/develop/contracts/STYLE.md).

Contributions are welcome! Please refer to
[Plugin's contributing guidelines](https://github.com/goplugin/pluginv3.0/blob/develop/docs/CONTRIBUTING.md) for detailed
contribution information.

Thank you!

### Changesets

We use [changesets](https://github.com/changesets/changesets) to manage versioning the contracts.

Every PR that modifies any configuration or code, should most likely accompanied by a changeset file.

To install `changesets`:
  1. Install `pnpm` if it is not already installed - [docs](https://pnpm.io/installation).
  2. Run `pnpm install`.

Either after or before you create a commit, run the `pnpm changeset` command in the `contracts` directory to create an accompanying changeset entry which will reflect on the CHANGELOG for the next release.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),

and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## License
Most of the contracts are licensed under the [MIT](https://choosealicense.com/licenses/mit/) license. 
An exception to this is the ccip folder, which defaults to be licensed under the [BUSL-1.1](./src/v0.8/ccip/LICENSE.md) license, however, there are a few exceptions

- `src/v0.8/ccip/applications/*` is licensed under the [MIT](./src/v0.8/ccip/LICENSE-MIT.md) license
- `src/v0.8/ccip/interfaces/*` is licensed under the [MIT](./src/v0.8/ccip/LICENSE-MIT.md) license
- `src/v0.8/ccip/libraries/{Client.sol, Internal.sol}` is licensed under the [MIT](./src/v0.8/ccip/LICENSE-MIT.md) license