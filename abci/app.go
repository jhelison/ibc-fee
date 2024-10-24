package abci

import (
	"crypto/sha256"

	"github.com/cometbft/cometbft/abci/types"
)

// App is the ABCI basic application with internal data
// This simulates a app and cuts down the a real consensus validation
type App struct {
	types.BaseApplication
	state              *State
	previousState      *State
	totalValidators    int64
	votesForBlock      int64
	consensusThreshold float64
}

// Assert the interface
var _ ABCIInterface = (*App)(nil)

// NewApp returns a new simulated app
func NewApp(totalValidators int64, consensusThreshold float64) *App {
	return &App{
		state:              &State{Height: 0, Data: make(map[string][]byte)},
		previousState:      &State{Height: 0, Data: make(map[string][]byte)},
		totalValidators:    totalValidators,
		consensusThreshold: consensusThreshold,
	}
}

// BeginBlock simulates the start of a block
// The main responsibility here is to start with a empty consensus state
func (app *App) BeginBlock(types.RequestBeginBlock) types.ResponseBeginBlock {
	// Prepare the app by updating the consensus state
	app.votesForBlock = 0

	// Copy the state as a backup
	app.previousState = &State{
		Height: app.state.Height,
		Data:   make(map[string][]byte),
	}
	for hash, tx := range app.state.Data {
		app.previousState.Data[hash] = tx
	}

	// We don't need to emit events though the simulation
	return types.ResponseBeginBlock{}
}

func (app *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	// Simplified processing of a TX
	txHash := HashTx(req.Tx)
	app.state.Data[txHash] = req.Tx

	// We can return a empty response for simplicity
	// But the event can be zero similarly to Cosmos-SDK
	return types.ResponseDeliverTx{Code: 0}
}

func (app *App) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	// Check the consensus
	if app.hasPassedConsensus() {
		// Increase the height
		app.state.Height += 1
		// If it has passed, we commit the state
		_ = app.Commit()
	} else {
		// If fail we rollback
		app.rollbackState()
	}

	// In this simulate we don't update the internal consensus state
	return types.ResponseEndBlock{}
}

func (app *App) rollbackState() {
	// Restore the old state
	// Normally this would be done by not committing the database changes
	app.state = app.previousState
}

// Vote is a helper function to vote in a block
func (app *App) Vote() {
	app.votesForBlock += 1
}

// hasPassedConsensus check if the votes were enough to pass consensus
func (app App) hasPassedConsensus() bool {
	votesPercent := (float64(app.votesForBlock) / float64(app.totalValidators))
	return votesPercent >= app.consensusThreshold
}

// GetState returns the current state, used on tests
func (app App) GetState() *State {
	return app.state
}

// hashTx hashes a TX with Sh256 to generate a key for the state
// Instead of using this we could use a strategy design pattern to make the hashing more flexible
func HashTx(tx []byte) string {
	// Create a new sha256 and write the TX
	h := sha256.New()
	h.Write(tx)

	// Return the stringified version of the tx
	return string(h.Sum(nil))
}
