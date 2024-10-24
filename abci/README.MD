# ABCI

This simulates a small implementation of the Application Blockchain interface:

- It implement a app with reduced ABCI interface and also reduces the logic behind TX processing

## Inner workings

The app has the following interface implemented:

- `BeginBlock`
  - Starts a new block, restart votes and initialize a state backup
- `DeliverTx`
  - Register a new TX on the block
  - No validations are done on the TX
- `EndBlock`
  - Ends a block and validates the consensus
  - If consensus was not reached, it rollback the state

## Files description

Description of each file and it's purpose:

- [App](./app.go)
  - This is the implementation of a simulated App with the ABCI interface
- [Interfaces](./interfaces.go)
  - The cut down ABCI interface
- [Types](./types.go)
  - A general state for the app with no implementation

Tests:

- Tests cover the following functionality:
  - Consensus pass
  - Consensus fail
  - Mixed consensus
- Tests can be found at:
  - [Tests](./app_test.go)
