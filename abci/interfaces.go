package abci

import (
	"github.com/cometbft/cometbft/abci/types"
)

// ABCIInterface is a cut down interface from: cometbft/abci/types/application.go
// This only has a basic structure with: BeginBlock, DeliverTx and EndBlock
type ABCIInterface interface {
	BeginBlock(types.RequestBeginBlock) types.ResponseBeginBlock
	DeliverTx(types.RequestDeliverTx) types.ResponseDeliverTx
	EndBlock(types.RequestEndBlock) types.ResponseEndBlock
	Commit() types.ResponseCommit
}
