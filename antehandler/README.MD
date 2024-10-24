# Antehandler

This is a general demonstration of the implementation of a fee antehandler:

- The fee antehandler takes the TX, validates the size, and charge coins based on the result
- The final formula is set as:

`Fee price * (tx bytes - Min Tx Size)`

## Inner workings

This antehandler imagines its implementation together with a feeHandler module:

- The module will have the params necessary to generate the new fee

## Files description

Description of each file and it's purpose:

- [The antehandler](./weighted_fee_ante.go)
  - This is the implementation of the new antehandler
- [Expected keepers](./weighted_fee_ante.go)
  - Definition of the interfaces used on the antehandler

Tests:

- Tests cover the following functionality:
  - TXs smaller than the minimum threshold
  - Txs bigger than the threshold
- Tests can be found at:
  - [Tests](./weighted_fee_ante_test.go)
