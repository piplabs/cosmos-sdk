package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetUbi returns the current distribution ubi.
func (k Keeper) GetUbi(ctx context.Context) (math.LegacyDec, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return params.Ubi, nil
}

// SetUbi sets new ubi
func (k Keeper) SetUbi(ctx context.Context, newUbi math.LegacyDec) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	params.Ubi = newUbi

	if err := params.Validate(); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return k.Params.Set(ctx, params)
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx context.Context) (enabled bool, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	return params.WithdrawAddrEnabled, nil
}
