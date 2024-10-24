# IBC

This simulates a simple broadcasting of packets between two apps:
- The simulations is a simplified version of what we have on chain though the IBC module

## Inner workings

The app has the following interface implemented:

- `SendPacket`
  - Creates and send the packet into a queue
  - The queue simulates a relayer
- `RecvPacket`
  - Receives the packet in the counterpart chain
  - And register it's data as a simulated TX
- `AcknowledgePacket`
  - Acknowledge that a packet was correct broadcast

## Files description

Description of each file and it's purpose:

- [App](./app.go)
  - This is the implementation of a simulated App with the IBC handler
- [Interfaces](./interfaces.go)
  - The cut down IBC ICS interface

Tests:

- Tests cover the following functionality:
  - Correct broadcast of packets
  - Broadcast of expired packets
  - Out of order packets
- Tests can be found at:
  - [Tests](./app_test.go)