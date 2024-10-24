package antehandler_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"ibc-fee/antehandler"
)

// Mock the FeeHandlerMock for tests
type FeeHandlerMock struct {
	params antehandler.FeeHandlerParams
}

var (
	// Definition of the default params for the mock module
	DefaultFeeBytePrice        = sdk.NewDecCoins(sdk.NewDecCoinFromDec("testcoin", sdk.OneDec()))
	DefaultMinTxSize    uint64 = 100 // Arbitrary size
)

// NewFeeHandlerMock returns a FeeHandlerMock
func NewFeeHandlerMock() FeeHandlerMock {
	return FeeHandlerMock{
		params: antehandler.FeeHandlerParams{
			FeeBytePrice: DefaultFeeBytePrice,
			MinTxSize:    DefaultMinTxSize,
		},
	}
}

// GetParams returns the FeeHandlerParams
func (fhm FeeHandlerMock) GetParams(ctx sdk.Context) antehandler.FeeHandlerParams {
	return fhm.params
}
