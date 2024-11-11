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
│   ├── v0.4
│   ├── v0.5
│   ├── v0.6
│   ├── v0.7
│   └── v0.8
└── abi # ABI json output
    ├── v0.4
    ├── v0.5
    ├── v0.6
    ├── v0.7
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

## License

[MIT](https://choosealicense.com/licenses/mit/)
