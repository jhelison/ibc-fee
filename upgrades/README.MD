# Upgrades

This is a general overview of Cosmos-SDK upgrades. This will focus on the application layer for upgrades.

## How an upgrade works and what is involved

An upgrade uses a few on-chain mechanisms:

1. Governance module

All starts with the governance module:

- A proposal must be submitted to allow the upgrade of a chain
- The proposal will have information about the upgrade, such as:
  - `Name`: The name of the upgrade, usually the target version
  - `Time`: A time for the upgrade can be set, although this is rarely used
  - `Height`: The height at the chain will halt and wait for an upgrade
  - `Info`: Extra information about the upgrade, such as links for the new binaries

2. Upgrade module

The upgrade module will hold the migrations and core logic for halting and resuming upgraded blockchains:

- The upgrade module will hold the migrations and the state changes for the new app version
- After the proposal passes the module will:
  - Halt the chain at the target block height
  - Apply the upgrade with state changes
  - Resume block production

## Step-by-step process

This is the step-by-step process of an upgrade:

1. Development

- The chain reaches a new version that requires an upgrade

2. Proposal submission

- A new proposal is published with the intent of upgrading the chain

3. Voting

- The proposal is voted and once passed the upgrade will be scheduled

4. Chain halting

- When the chain reaches the target height, all nodes halt and start to wait for the upgrade

5. Node upgrade

- The operator installs the new binary
- The process can be automated if using Cosmosvisor

6. State Migration

- When initialized the planned migrations will be applied

7. Chain resume

- After enough validators have done the upgrade procedure, the chain resumes block production

## Risks and Mitigations

The following are risks and possible mitigations during an upgrade:

**Failed Upgrade**

- An error on the migration script happens, causing an upgrade failure

**Extended Downtime**

- Validators take a long time to complete the upgrade resulting in a prolonged downtime

Mitigations:

**Devnets and Testnets**

- Upgrades can be tested on Devnets and Testnets
- This reduces the possibility of a failed migration during a mainnet upgrade

**Planning**

- Allow the upgrade height to land at a good time for validators
- This increases the chance of a fast and successful upgrade

## Mock Upgrade Handler Implementation

This is a simple mock upgrade handler implementation:

### Registering the upgrade handlers

```go
// The upgrade handler sample
func (app App) RegisterUpgradeHandlers() {
	// if there is an override list, use that instead, for integration tests
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				app.configurator,
				&app.AppKeepers,
			),
		)
	}
}
```

### Creating an upgrade handler

```go
// CreateUpgradeHandler defines a new upgrade that migrates the db to a new version
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, newVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, newVM)
		if err != nil {
			return vm, err
		}

    // Here you can apply any state migration and iterate with modules

		return newVM, err
	}
}
```
