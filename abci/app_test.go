package abci_test

import (
	"testing"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/assert"

	"ibc-fee/abci"
)

// TestConsensusPass tests a perfect consensus on the simulated app
func TestConsensusPass(t *testing.T) {
	app := abci.NewApp(4, 0.66)

	// Start a new block
	app.BeginBlock(types.RequestBeginBlock{})

	// Deliver TXs
	app.DeliverTx(types.RequestDeliverTx{Tx: []byte("tx1")})
	app.DeliverTx(types.RequestDeliverTx{Tx: []byte("tx2")})

	// Simulate validators vote on the block
	voteOnApp(4, app)

	// Finalize the block
	app.EndBlock(types.RequestEndBlock{})

	// Assert the new height
	currentState := app.GetState()
	assert.Equal(t, int64(1), currentState.Height)

	// Assert that TXs exists
	tx1Hash := abci.HashTx([]byte("tx1"))
	assert.Equal(t, []byte("tx1"), currentState.Data[tx1Hash])
	tx2Hash := abci.HashTx([]byte("tx2"))
	assert.Equal(t, []byte("tx2"), currentState.Data[tx2Hash])
}

// TestFailedConsensus test a failed consensus by not reaching the vote threshold
func TestFailedConsensus(t *testing.T) {
	app := abci.NewApp(4, 0.66)

	// Start a new block
	app.BeginBlock(types.RequestBeginBlock{})

	// Deliver TXs
	app.DeliverTx(types.RequestDeliverTx{Tx: []byte("tx1")})
	app.DeliverTx(types.RequestDeliverTx{Tx: []byte("tx2")})

	// Simulate vote, but only 50%
	voteOnApp(2, app)

	// Finalize the block
	app.EndBlock(types.RequestEndBlock{})

	// Assert that height remains the same
	currentState := app.GetState()
	assert.Equal(t, int64(0), currentState.Height)

	// The list of TXs must be empty
	assert.Empty(t, currentState.Data)
}

// TestMixedConsensus test a passing consensus and them a failed consensus
// The first set of messages should be commit but not the failed set
func TestMixedConsensus(t *testing.T) {
	app := abci.NewApp(10, 0.50) // Use different values just to cover more field

	// PASSING BLOCK 0

	// Start a new block
	app.BeginBlock(types.RequestBeginBlock{})

	// Deliver TXs
	app.DeliverTx(types.RequestDeliverTx{Tx: []byte("passing_tx")})

	// Simulate vote, 50% should be on the threshold and should pass
	voteOnApp(5, app)

	// Finalize the block
	app.EndBlock(types.RequestEndBlock{})

	// FAILED BLOCK 1

	// Start a new block
	app.BeginBlock(types.RequestBeginBlock{})

	// Deliver TXs
	app.DeliverTx(types.RequestDeliverTx{Tx: []byte("failed_tx")})

	// Simulate vote, but just bellow consensus
	voteOnApp(4, app)

	// Finalize the block
	app.EndBlock(types.RequestEndBlock{})

	// ASSERT

	// We should be on block 1
	currentState := app.GetState()
	assert.Equal(t, int64(1), currentState.Height)

	// Only one TX must exist and should be the passed one
	assert.Equal(t, 1, len(currentState.Data))
	passingTxHash := abci.HashTx([]byte("passing_tx"))
	assert.Equal(t, []byte("passing_tx"), currentState.Data[passingTxHash])
}

// voteOnApp is a helper function to test the voting process on the app
func voteOnApp(totalVotes int, app *abci.App) {
	for range totalVotes {
		app.Vote()
	}
}
