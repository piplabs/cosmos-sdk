package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// zero fee pool
func InitialFeePool() FeePool {
	return FeePool{
		Ubi: sdk.DecCoins{},
	}
}

// ValidateGenesis validates the fee pool for a genesis state
func (f FeePool) ValidateGenesis() error {
	if f.Ubi.IsAnyNegative() {
		return fmt.Errorf("negative ubi in distribution fee pool, is %v",
			f.Ubi)
	}

	return nil
}
