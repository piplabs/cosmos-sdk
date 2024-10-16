package types

import (
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
)

func NewPeriodDelegation(
	periodDelegationID string, shares math.LegacyDec, rewardsShares math.LegacyDec, periodType PeriodType, endTime time.Time,
) PeriodDelegation {
	return PeriodDelegation{
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
