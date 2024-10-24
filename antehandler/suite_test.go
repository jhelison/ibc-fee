// Suite for the weighted antehandler test
// This easy out the testing by implementing modules as mocks
// This implementation was heavily expired from cosmos-sdk/x/auth/ante/testutil_test.go
package antehandler_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtestutil "github.com/cosmos/cosmos-sdk/x/auth/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AnteTestSuite is a test suite to be used on the weighted fee antehandler tests
type AnteTestSuite struct {
	ctx        sdk.Context
	bankKeeper *authtestutil.MockBankKeeper
	feeHandler FeeHandlerMock
}

// SetupTest setups a new test with mock bank implementation
func SetupTestSuite(t *testing.T, isCheckTx bool) *AnteTestSuite {
	// Initialize the mock bank and controller
	suite := &AnteTestSuite{}
	ctrl := gomock.NewController(t)
	suite.bankKeeper = authtestutil.NewMockBankKeeper(ctrl)

	// Initialize a new Key value store and a testing context
	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	suite.ctx = testCtx.Ctx.WithIsCheckTx(isCheckTx).WithBlockHeight(1)

	// Initialize the feeHandler mock
	suite.feeHandler = NewFeeHandlerMock()

	return suite
}
