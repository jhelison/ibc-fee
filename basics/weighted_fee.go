package basics

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

const ErrForbiddenTXTypeFormat = "error parsing tx into it's proto reference"

// WeightedFeeDecorator is the decorator responsible of charging extra fees based on a TX size
type WeightedFeeDecorator struct {
	bankKeeper BankKeeper
}

// NewWeightedFeeDecorator returns a new weighted fee decorator
func NewWeightedFeeDecorator(bk BankKeeper) WeightedFeeDecorator {
	// Returns the object
	return WeightedFeeDecorator{
		bankKeeper: bk,
	}
}

// AnteHandler executes the effective antehandler function
// If checks if the the tx
func (wfd WeightedFeeDecorator) AnteHandler(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Encode the TX into the proto reference
	protoTx, ok := tx.(*txtypes.Tx)
	if !ok {
		return ctx, errorsmod.Wrapf(
			errortypes.ErrInvalidType,
			ErrForbiddenTXTypeFormat,
		)
	}

	// Marshal the TX into bytes
	txBytes, err := protoTx.Marshal()
	if err != nil {
		return ctx, err
	}

	// Calculate the new ratio
	fmt.Println(txBytes)

	// Continue the decorator execution with the next function
	return next(ctx, tx, simulate)
}
