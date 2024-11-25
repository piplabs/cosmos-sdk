package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func (k Keeper) GetUbiBalanceByDenom(ctx context.Context, denom string) (math.Int, error) {
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return math.Int{}, err
	}

	return feePool.Ubi.AmountOf(denom).TruncateInt(), nil
}

func (k Keeper) WithdrawUbiByDenomToModule(ctx context.Context, denom string, recipientModule string) (sdk.Coin, error) {
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	coin := sdk.NewCoin(denom, feePool.Ubi.AmountOf(denom).TruncateInt())
	coins := sdk.NewCoins(coin)

	// NOTE the ubi pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the ubi pool
	// must be reduced separately from the SendCoinsFromModuleToModule call
	newPool, negative := feePool.Ubi.SafeSub(sdk.NewDecCoinsFromCoins(coins...))
	if negative {
		return sdk.Coin{}, types.ErrBadDistribution
	}
	feePool.Ubi = newPool

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientModule, coins); err != nil {
		return sdk.Coin{}, err
	}

	if err := k.FeePool.Set(ctx, feePool); err != nil {
		return sdk.Coin{}, err
	}

	return coin, nil
}
