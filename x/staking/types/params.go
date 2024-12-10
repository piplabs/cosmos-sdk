package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7 * 3

	// Default maximum number of bonded validators
	DefaultMaxValidators uint32 = 100

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint32 = 7

	// DefaultHistorical entries is 10000. Apps that don't use IBC can ignore this
	// value by not adding the staking module to the application module manager's
	// SetOrderBeginBlockers.
	DefaultHistoricalEntries uint32 = 10000

	DefaultFlexiblePeriodType = 0

	DefaultLockedTokenType   = 0
	DefaultSingularityHeight = 403200 // 14 days with 3 seconds block time
)

// DefaultMinCommissionRate is set to 0%
var DefaultMinCommissionRate = math.LegacyZeroDec()

var DefaultMinDelegation = math.NewInt(1)

var DefaultPeriods = []Period{
	{
		PeriodType:        0,
		Duration:          time.Duration(0),    // flexible
		RewardsMultiplier: math.LegacyOneDec(), // 1
	},
	{
		PeriodType:        1,
		Duration:          time.Hour * 24 * 90,                // 90 days
		RewardsMultiplier: math.LegacyNewDecWithPrec(1051, 3), // 1.051
	},
	{
		PeriodType:        2,
		Duration:          time.Hour * 24 * 360,              // 360 days
		RewardsMultiplier: math.LegacyNewDecWithPrec(116, 2), // 1.16
	},
	{
		PeriodType:        3,
		Duration:          time.Hour * 24 * 540,              // 540 days
		RewardsMultiplier: math.LegacyNewDecWithPrec(134, 2), // 1.34
	},
}

var DefaultTokenTypes = []TokenTypeInfo{
	{
		TokenType:         0,                               // Locked
		RewardsMultiplier: math.LegacyNewDecWithPrec(5, 1), // 0.5
	},
	{
		TokenType:         1,                   // Unlocked
		RewardsMultiplier: math.LegacyOneDec(), // 1
	},
}

// NewParams creates a new Params instance
func NewParams(
	unbondingTime time.Duration, maxValidators, maxEntries, historicalEntries uint32, bondDenom string, minCommissionRate math.LegacyDec,
	minDelegation math.Int, flexiblePeriodType int32, periods []Period, lockedTokenType int32, tokenTypes []TokenTypeInfo, singulariyHeight uint64,
) Params {
	return Params{
		UnbondingTime:      unbondingTime,
		MaxValidators:      maxValidators,
		MaxEntries:         maxEntries,
		HistoricalEntries:  historicalEntries,
		BondDenom:          bondDenom,
		MinCommissionRate:  minCommissionRate,
		MinDelegation:      minDelegation,
		FlexiblePeriodType: flexiblePeriodType,
		Periods:            periods,
		LockedTokenType:    lockedTokenType,
		TokenTypes:         tokenTypes,
		SingularityHeight:  singulariyHeight,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMaxValidators,
		DefaultMaxEntries,
		DefaultHistoricalEntries,
		sdk.DefaultBondDenom,
		DefaultMinCommissionRate,
		DefaultMinDelegation,
		DefaultFlexiblePeriodType,
		DefaultPeriods,
		DefaultLockedTokenType,
		DefaultTokenTypes,
		DefaultSingularityHeight,
	)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	err = cdc.Unmarshal(value, &params)
	if err != nil {
		return
	}

	return
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}

	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}

	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}

	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}

	if err := validateMinCommissionRate(p.MinCommissionRate); err != nil {
		return err
	}

	if err := validateHistoricalEntries(p.HistoricalEntries); err != nil {
		return err
	}

	if err := validateFlexiblePeriodType(p.FlexiblePeriodType); err != nil {
		return err
	}

	if err := validatePeriods(p.Periods); err != nil {
		return err
	}

	if err := validateLockedTokenType(p.LockedTokenType); err != nil {
		return err
	}

	if err := validateTokenTypes(p.TokenTypes); err != nil {
		return err
	}

	if err := validateSingularityHeight(p.SingularityHeight); err != nil {
		return err
	}

	return nil
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMaxValidators(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidatePowerReduction(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(math.NewInt(1)) {
		return fmt.Errorf("power reduction cannot be lower than 1")
	}

	return nil
}

func validateMinCommissionRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("minimum commission rate cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("minimum commission rate cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("minimum commission rate cannot be greater than 100%%: %s", v)
	}

	return nil
}

func validateFlexiblePeriodType(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid flexible period type: %d", v)
	}

	return nil
}

func validatePeriods(i interface{}) error {
	periods, ok := i.([]Period)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, period := range periods {
		if period.PeriodType < 0 {
			return fmt.Errorf("invalid period type: %d", period.PeriodType)
		}
		if !period.RewardsMultiplier.IsPositive() {
			return fmt.Errorf("invalid period rewards multiplier: %s", period.RewardsMultiplier.String())
		}
	}

	return nil
}

func validateLockedTokenType(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid locked token type: %d", v)
	}

	return nil
}

func validateTokenTypes(i interface{}) error {
	tokenTypes, ok := i.([]TokenTypeInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, tokenType := range tokenTypes {
		if tokenType.TokenType < 0 {
			return fmt.Errorf("invalid token type: %d", tokenType.TokenType)
		}
		if !tokenType.RewardsMultiplier.IsPositive() {
			return fmt.Errorf("invalid token rewards multiplier: %s", tokenType.RewardsMultiplier.String())
		}
	}

	return nil
}

func validateSingularityHeight(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("singularity height must be positive: %d", v)
	}

	return nil
}
