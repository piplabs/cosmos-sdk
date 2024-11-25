package keeper

import (
	"context"
	"errors"
	"time"

	storetypes "cosmossdk.io/store/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetOrCreatePeriodDelegation gets the period delegation or creates a new one.
func (k Keeper) GetOrCreatePeriodDelegation(
	ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress,
	periodDelegationID string,
	periodType int32, endTime time.Time,
) (types.PeriodDelegation, error) {
	periodDelegation, err := k.GetPeriodDelegation(ctx, delAddr, valAddr, periodDelegationID)
	if errors.Is(err, types.ErrNoPeriodDelegation) {
		delAddrStr, err := k.authKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return types.PeriodDelegation{}, err
		}
		valAddrStr, err := k.validatorAddressCodec.BytesToString(valAddr)
		if err != nil {
			return types.PeriodDelegation{}, err
		}
		periodDelegation = types.NewPeriodDelegation(
			delAddrStr, valAddrStr,
			periodDelegationID, math.LegacyZeroDec(), math.LegacyZeroDec(), periodType, endTime,
		)
	} else if err != nil {
		return types.PeriodDelegation{}, err
	}

	return periodDelegation, nil
}

// GetAllPeriodDelegationsByDelAndValAddr returns all period delegations by delAddr and valAddr.
func (k Keeper) GetAllPeriodDelegationsByDelAndValAddr(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) ([]types.PeriodDelegation, error) {
	store := k.storeService.OpenKVStore(ctx)

	periodDelegationsKey := types.GetPeriodDelegationsKey(delAddr, valAddr)
	iterator, err := store.Iterator(periodDelegationsKey, storetypes.PrefixEndBytes(periodDelegationsKey))
	if err != nil {
		return nil, err
	}

	periodDelegations := make([]types.PeriodDelegation, 0)
	for ; iterator.Valid(); iterator.Next() {
		periodDelegation, err := types.UnmarshalPeriodDelegation(k.cdc, iterator.Value())
		if err != nil {
			return nil, err
		}
		periodDelegations = append(periodDelegations, periodDelegation)
	}

	return periodDelegations, nil
}

// GetAllPeriodDelegations returns all period delegations
func (k Keeper) GetAllPeriodDelegations(ctx context.Context) (periodDelegations []types.PeriodDelegation, err error) {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.PeriodDelegationKey, storetypes.PrefixEndBytes(types.PeriodDelegationKey))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		periodDelegation, err := types.UnmarshalPeriodDelegation(k.cdc, iterator.Value())
		if err != nil {
			return nil, err
		}
		periodDelegations = append(periodDelegations, periodDelegation)
	}

	return
}

// GetPeriodDelegation returns a specific period delegation.
func (k Keeper) GetPeriodDelegation(
	ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, periodDelegationID string,
) (types.PeriodDelegation, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetPeriodDelegationKey(delAddr, valAddr, periodDelegationID)

	value, err := store.Get(key)
	if err != nil {
		return types.PeriodDelegation{}, err
	} else if value == nil {
		return types.PeriodDelegation{}, types.ErrNoPeriodDelegation
	}

	return types.UnmarshalPeriodDelegation(k.cdc, value)
}

// SetPeriodDelegation sets a period delegation.
func (k Keeper) SetPeriodDelegation(
	ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, periodDelegation types.PeriodDelegation,
) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(
		types.GetPeriodDelegationKey(delAddr, valAddr, periodDelegation.PeriodDelegationId),
		types.MustMarshalPeriodDelegation(k.cdc, periodDelegation),
	)
}

// RemovePeriodDelegation removes a period delegation
func (k Keeper) RemovePeriodDelegation(
	ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, periodDelegationID string,
) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.GetPeriodDelegationKey(delAddr, valAddr, periodDelegationID))
}
