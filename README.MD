# IBC / ABCI / Antehandler

This repository consists of multiple experiments and simulations related to Cosmos-SDK. It covers:
- Antehandler
- IBC packet broadcasting
- ABCI interface implementation
- Description of upgrades

## Projects

Each project can be found in it's own directory, but summing up, each contains:
- The implementation
- Tests
- A readme with a description of the functionality

This is a general description of each:
- [Antehandler](./antehandler/README.MD)
  - Implements a antehandler capable of charging fees based on a threshold and TX size
- [ABCI](./abci/README.MD)
  - Implements the ABCI interface into a simulated app
- [IBC](./ibc/README.MD)
  - Implements packet handling into a simulated app
- [Upgrades](./upgrades/README.MD)
  - Has the general description of upgrades

## How to run this project

Since this project is mostly composed of POCs, it only has tests.
You can run the tests with:

```bash
make test-unit
```