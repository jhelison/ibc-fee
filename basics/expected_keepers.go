package basics

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the interface of the banking Keeper used on the weighted_fee ante handler
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
