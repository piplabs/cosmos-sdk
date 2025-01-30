package types

import (
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FlexiblePeriodDelegationID = "0"
)

// PeriodDelegations is a collection of period delegations
type PeriodDelegations []PeriodDelegation

func NewPeriodDelegation(
	delegatorAddr, validatorAddr, periodDelegationID string, shares, rewardsShares math.LegacyDec,
	periodType int32, endTime time.Time,
) PeriodDelegation {
	return PeriodDelegation{
		DelegatorAddress:   delegatorAddr,
		ValidatorAddress:   validatorAddr,
		PeriodDelegationId: periodDelegationID,
		Shares:             shares,
		RewardsShares:      rewardsShares,
		PeriodType:         periodType,
		EndTime:            endTime,
	}
}

// MustMarshalPeriodDelegation returns the period delegation bytes. Panics if fails
func MustMarshalPeriodDelegation(cdc codec.BinaryCodec, periodDelegation PeriodDelegation) []byte {
	return cdc.MustMarshal(&periodDelegation)
}

// MustUnmarshalPeriodDelegation return the unmarshaled period delegation from bytes.
// Panics if fails.
func MustUnmarshalPeriodDelegation(cdc codec.BinaryCodec, value []byte) PeriodDelegation {
	periodDelegation, err := UnmarshalPeriodDelegation(cdc, value)
	if err != nil {
		panic(err)
	}

	return periodDelegation
}

// return the period delegation
func UnmarshalPeriodDelegation(cdc codec.BinaryCodec, value []byte) (periodDelegation PeriodDelegation, err error) {
	err = cdc.Unmarshal(value, &periodDelegation)
	return periodDelegation, err
}

// Client Types

// PeriodDelegationResponses is a collection of PeriodDelegationResponse
type PeriodDelegationResponses []PeriodDelegationResponse

// NewPeriodDelegationResp creates a new PeriodDelegationResponse instance
func NewPeriodDelegationResp(
	periodDelegation PeriodDelegation, balance sdk.Coin,
) PeriodDelegationResponse {
	return PeriodDelegationResponse{
		PeriodDelegation: periodDelegation,
		Balance:          balance,
	}
}
