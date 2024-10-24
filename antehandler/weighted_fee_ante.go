// Antehandler to charge fees based on the tx bytes size
// This implementation uses a simulated module called FeeHandler
// The simulated module stores information such as fee prices per byte and minimum fee size to charge fees
// This module is inspired on Cosmos-SDK Fee antehandler, but with simplifications:
// - Fee grants are disabled
// - We consider the FeeCollector acc set
package antehandler

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Error messages
const (
	ErrFeeTxDecode = "error parsing tx into FeeTx"

	// Events to be emitted since we don't have a module behind
	AttributeKeyBytesFee = "bytes_fee"
)

// Assert that the AnteDecorator function is really being implemented
var _ sdk.AnteDecorator = (*WeightedFeeDecorator)(nil)

// WeightedFeeDecorator is the decorator responsible of charging extra fees based on a TX size
type WeightedFeeDecorator struct {
	bankKeeper BankKeeper
	feeHandler FeeHandler
}

// NewWeightedFeeDecorator returns a new weighted fee decorator
func NewWeightedFeeDecorator(bk BankKeeper, fh FeeHandler) WeightedFeeDecorator {
	// Returns the object
	return WeightedFeeDecorator{
		bankKeeper: bk,
		feeHandler: fh,
	}
}

// AnteHandle executes the effective antehandler function
// It charges fees on top of normal fees based on the feeHandler module
// Only extra bytes are charged from the user
// Fees are calculated as:
// Fee price * (tx bytes - Min Tx Size)
func (wfd WeightedFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Pass the call to the check and deduct fee
	err := wfd.checkDeductFee(ctx, tx)
	if err != nil {
		return ctx, err
	}

	// Continue the decorator execution with the next function
	return next(ctx, tx, simulate)
}

// checkDeductFee checks the tx and deducts the fee
func (wfd WeightedFeeDecorator) checkDeductFee(ctx sdk.Context, tx sdk.Tx) error {
	// Validate if the Tx has the minimum size to get the extra fees
	// In a production environment we should use the chain encoders though the set codecs
	txBytes, err := authtx.DefaultTxEncoder()(tx)
	if err != nil {
		return err
	}

	// Get the feeHandler params
	feeHandlerParams := wfd.feeHandler.GetParams(ctx)

	// Check if our TX will pay extra fees
	if len(txBytes) <= int(feeHandlerParams.MinTxSize) {
		return nil
	}

	// Parse the TX as a FeeTx
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return errorsmod.Wrap(
			errortypes.ErrTxDecode, ErrFeeTxDecode,
		)
	}

	// Get the fee payer from the TX
	feePayer := feeTx.FeePayer()

	// Calculate the total fee, but only for the additional bytes
	extraBytes := len(txBytes) - int(feeHandlerParams.MinTxSize)
	totalFee := calculateFeeForBytes(int64(extraBytes), feeHandlerParams.FeeBytePrice)

	// Charge the extra fee from the user
	if !totalFee.IsZero() {
		err := wfd.bankKeeper.SendCoinsFromAccountToModule(ctx, feePayer, types.FeeCollectorName, totalFee)
		if err != nil {
			return err
		}
	}

	// Emit events
	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(AttributeKeyBytesFee, totalFee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, feePayer.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	// No errors were reached until now
	return nil
}

// calculateFeeForBytes calculate the fees for a txbytes
// It use the formula: Fee price * size
// This also truncate the decimals
func calculateFeeForBytes(size int64, bytesPrice sdk.DecCoins) sdk.Coins {
	bytesValue := bytesPrice.MulDec(sdk.NewDec(size))
	total, _ := bytesValue.TruncateDecimal()
	return total
}
