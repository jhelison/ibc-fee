package antehandler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeHandlerParams are the params for the simulated feeHandlerModule
type FeeHandlerParams struct {
	// FeeBytePrice is the price for each byte in a TX
	FeeBytePrice sdk.DecCoins
	// The minimum size a TX must have
	MinTxSize uint64
}

// BankKeeper defines the interface of the banking Keeper used on the weighted_fee ante handler
type BankKeeper interface {
	IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error
	SendCoins(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// FeeHandler defines a simulated feeHandler module
type FeeHandler interface {
	GetParams(ctx sdk.Context) FeeHandlerParams
}
