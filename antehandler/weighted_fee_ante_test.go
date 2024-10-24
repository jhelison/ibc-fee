package antehandler_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktyppes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	ante "ibc-fee/antehandler"
)

// TestWeightedFeeAnte tests the weighted fee antehandler
func TestWeightedFeeAnte(t *testing.T) {
	// Prepare the testing data
	accAddr1 := sdk.AccAddress([]byte("acc1"))
	accAddr2 := sdk.AccAddress([]byte("acc2"))

	// All the test cases
	testCases := []struct {
		name        string
		msgs        []sdk.Msg
		expectedFee sdk.Coins
	}{
		{
			name: "No fee, single bank send message",
			msgs: []sdk.Msg{
				banktyppes.NewMsgSend( // The fee handler has been set to have exactly this TX as limit
					accAddr1,
					accAddr2,
					sdk.NewCoins(sdk.NewCoin("utestcoin", sdk.OneInt())),
				),
			},
			expectedFee: nil,
		},
		{
			name: "Fee, two bank send messages",
			msgs: []sdk.Msg{
				banktyppes.NewMsgSend(
					accAddr1,
					accAddr2,
					sdk.NewCoins(sdk.NewCoin("utestcoin", sdk.OneInt())),
				),
				banktyppes.NewMsgSend(
					accAddr1,
					accAddr2,
					sdk.NewCoins(sdk.NewCoin("utestcoin", sdk.OneInt())),
				),
			},
			expectedFee: sdk.NewCoins(sdk.NewCoin("testcoin", math.NewInt(95))), // This is the fee set on antehandler/mock_test.go. Each Bank msg is 95 above the limit
		},
		{
			name: "No fee, ibc message",
			msgs: []sdk.Msg{
				&ibcclienttypes.MsgCreateClient{
					Signer: accAddr1.String(),
				},
			},
			expectedFee: nil,
		},
		{
			name: "Fee, bank + ibc message",
			msgs: []sdk.Msg{
				banktyppes.NewMsgSend(
					accAddr1,
					accAddr2,
					sdk.NewCoins(sdk.NewCoin("utestcoin", sdk.OneInt())),
				),
				&ibcclienttypes.MsgCreateClient{
					Signer: accAddr1.String(),
				},
			},
			expectedFee: sdk.NewCoins(sdk.NewCoin("testcoin", math.NewInt(64))), // IBC message has a size of 64
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Case %s", tc.name), func(t *testing.T) {
			// At each run we restart our setup with a new fee decorator
			s := SetupTestSuite(t, false)
			dfd := ante.NewWeightedFeeDecorator(s.bankKeeper, s.feeHandler)

			// We initialize a new chain ante decorator with terminator
			antehandler := sdk.ChainAnteDecorators(dfd)

			// Build a new TX based on the passed MSGs
			tx := createTX(t, tc.msgs)

			// If the Tx is expected to have a new fee, we prepare the mockBank to handle the call
			if tc.expectedFee != nil {
				// Expect the call with the correct balance, this is our main validation on the total fees
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), tc.expectedFee).Return(sdkerrors.ErrInsufficientFunds)
			}

			// Run the antehandler
			_, err := antehandler(s.ctx, tx, false)

			// If a expected fee has been passed we will get insufficient funds error
			if tc.expectedFee != nil {
				require.ErrorIs(t, err, sdkerrors.ErrInsufficientFunds)
			} else {
				require.NoError(t, err)
			}

		})
	}
}

// createTX creates a new testing tx from current encoding
func createTX(t *testing.T, msgs []sdk.Msg) signing.Tx {
	// Create the TX
	encodingConfig := testutil.MakeTestEncodingConfig()
	txBuilder := encodingConfig.TxConfig.NewTxBuilder()

	// Set the msgs
	err := txBuilder.SetMsgs(msgs...)
	require.NoError(t, err)

	return txBuilder.GetTx()
}
