package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
)

// DefaultParams returns default distribution parameters
func DefaultParams() Params {
	return Params{
		Ubi:                 math.LegacyNewDecWithPrec(2, 2), // 2%
		BaseProposerReward:  math.LegacyZeroDec(),            // deprecated
		BonusProposerReward: math.LegacyZeroDec(),            // deprecated
		WithdrawAddrEnabled: true,
		MaxUbi:              math.LegacyNewDecWithPrec(2, 1), // 20%
	}
}

func (p Params) Validate() error {
	if err := validateUbi(p.Ubi); err != nil {
		return err
	}

	if err := validateMaxUbi(p.MaxUbi); err != nil {
		return err
	}

	if p.Ubi.GT(p.MaxUbi) {
		return errors.New("ubi must be less than or equal to max ubi")
	}

	return nil
}

// ValidateBasic performs basic validation on distribution parameters.
func (p Params) ValidateBasic() error {
	if err := validateUbi(p.Ubi); err != nil {
		return err
	}

	return validateMaxUbi(p.MaxUbi)
}

func validateUbi(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("ubi must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("ubi must be positive: %s", v)
	}
	if v.GTE(math.LegacyOneDec()) {
		return fmt.Errorf("ubi too large: %s", v)
	}

	return nil
}

func validateMaxUbi(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("max ubi must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("max ubi must be positive: %s", v)
	}
	if v.GTE(math.LegacyOneDec()) {
		return fmt.Errorf("max ubi too large: %s", v)
	}

	return nil
}

func validateWithdrawAddrEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
