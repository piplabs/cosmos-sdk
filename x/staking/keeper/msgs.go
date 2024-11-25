package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k msgServer) ValidateCreateValidatorMsg(ctx context.Context, msg *types.MsgCreateValidator) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	valAddr, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	// skip validation that if the token is greater than minimum self delegation amount at genesis block
	skipMinSelfDelValidation := sdkCtx.BlockHeight() == 0
	if err := msg.Validate(k.validatorAddressCodec, skipMinSelfDelValidation); err != nil {
		return err
	}

	minCommRate, err := k.MinCommissionRate(ctx)
	if err != nil {
		return err
	}
	if msg.Commission.Rate.LT(minCommRate) {
		return errorsmod.Wrapf(types.ErrCommissionLTMinRate, "cannot set validator commission to less than minimum rate of %s", minCommRate)
	}

	// check to see if the pubkey or sender has been registered before
	if _, err := k.GetValidator(ctx, valAddr); err == nil {
		return types.ErrValidatorOwnerExists
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, err := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); err == nil {
		return types.ErrValidatorPubKeyExists
	}

	bondDenom, err := k.BondDenom(ctx)
	if err != nil {
		return err
	}

	if msg.Value.Denom != bondDenom {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Value.Denom, bondDenom,
		)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}

	cp := sdkCtx.ConsensusParams()
	if cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return errorsmod.Wrapf(
				types.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	minDelegation, err := k.MinDelegation(ctx)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get min delegation")
	}
	// minimum self delegation should be greater than or equal to chain-side min delegation
	if msg.MinSelfDelegation.LT(minDelegation) {
		return types.ErrMinSelfDelegationBelowMinDelegation
	}
	// self delegation amount must be greater than or equal to minimum self delegation when creating validator
	if msg.Value.Amount.LT(msg.MinSelfDelegation) && !skipMinSelfDelValidation {
		return types.ErrSelfDelegationBelowMinimum
	}

	// validate token type that the validator is supporting
	if _, err := k.GetTokenTypeInfo(ctx, msg.SupportTokenType); err != nil {
		return err
	}

	// self delegation is flexible period delegation
	flexiblePeriodType, err := k.GetFlexiblePeriodType(ctx)
	if err != nil {
		return err
	}
	// check period type info
	if _, err := k.GetPeriodInfo(ctx, flexiblePeriodType); err != nil {
		return err
	}

	return nil
}

func (k msgServer) ValidateEditValidatorMsg(ctx context.Context, msg *types.MsgEditValidator) error {
	valAddr, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	// validator must already be registered
	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	if msg.Description == (types.Description{}) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil {
		if !msg.MinSelfDelegation.IsPositive() {
			return errorsmod.Wrap(
				sdkerrors.ErrInvalidRequest,
				"minimum self delegation must be a positive integer",
			)
		}

		if !msg.MinSelfDelegation.GTE(validator.MinSelfDelegation) {
			return types.ErrMinSelfDelegationDecreased
		}

		if msg.MinSelfDelegation.GT(validator.Tokens) {
			return types.ErrSelfDelegationBelowMinimum
		}
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(math.LegacyOneDec()) || msg.CommissionRate.IsNegative() {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}

		minCommissionRate, err := k.MinCommissionRate(ctx)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}

		if msg.CommissionRate.LT(minCommissionRate) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "commission rate cannot be less than the min commission rate %s", minCommissionRate.String())
		}
	}

	return nil
}

func (k msgServer) ValidateDelegateMsg(ctx context.Context, msg *types.MsgDelegate) error {
	valAddr, valErr := k.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if valErr != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", valErr)
	}

	if _, err := k.authKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	minDelegation, err := k.MinDelegation(ctx)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get min delegation")
	}
	// delegation amount should be greater than or equal to chain-side min delegation
	if msg.Amount.Amount.LT(minDelegation) {
		return types.ErrDelegationBelowMinimum
	}

	bondDenom, err := k.BondDenom(ctx)
	if err != nil {
		return err
	}

	if msg.Amount.Denom != bondDenom {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	// check if the validator exists
	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	// validate token type info
	if _, err := k.GetTokenTypeInfo(ctx, validator.GetSupportTokenType()); err != nil {
		return err
	}
	// validate period info
	if _, err := k.GetPeriodInfo(ctx, msg.PeriodType); err != nil {
		return err
	}
	// validate period type and delegation id
	flexiblePeriodType, err := k.GetFlexiblePeriodType(ctx)
	if err != nil {
		return err
	}
	if (msg.PeriodType == flexiblePeriodType) != (msg.PeriodDelegationId == types.FlexiblePeriodDelegationID) {
		return types.ErrPeriodDelegationIDMismatch
	}

	return nil
}

func (k msgServer) ValidateUndelegateMsg(ctx context.Context, msg *types.MsgUndelegate) error {
	validatorAddress, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	minUndelegation, err := k.MinDelegation(ctx)
	if err != nil {
		return err
	}
	// undelegation amount must be greater than or equal to minimum undelegation
	if msg.Amount.Amount.LT(minUndelegation) {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"undelegation amount is less than the minimum undelegation amount",
		)
	}

	bondDenom, err := k.BondDenom(ctx)
	if err != nil {
		return err
	}

	if msg.Amount.Denom != bondDenom {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	// check validator existence
	if _, err := k.GetValidator(ctx, validatorAddress); err != nil {
		return err
	}
	// check delegation existence
	if _, err := k.GetDelegation(ctx, delegatorAddress, validatorAddress); err != nil {
		return err
	}
	// check period delegation existence
	if _, err := k.GetPeriodDelegation(ctx, delegatorAddress, validatorAddress, msg.PeriodDelegationId); err != nil {
		return err
	}

	return nil
}

func (k msgServer) ValidateBeginRedelegateMsg(ctx context.Context, msg *types.MsgBeginRedelegate) error {
	if msg.ValidatorSrcAddress == msg.ValidatorDstAddress {
		return types.ErrSelfRedelegation
	}

	valSrcAddr, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorSrcAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}

	valDstAddr, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorDstAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}

	delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	srcVal, err := k.GetValidator(ctx, valSrcAddr)
	if err != nil {
		return err
	}

	dstVal, err := k.GetValidator(ctx, valDstAddr)
	if err != nil {
		return err
	}

	if srcVal.SupportTokenType != dstVal.SupportTokenType {
		return types.ErrTokenTypeMismatch
	}

	if err := k.ValidateUndelegateMsg(ctx, &types.MsgUndelegate{
		ValidatorAddress:   msg.ValidatorSrcAddress,
		DelegatorAddress:   msg.DelegatorAddress,
		Amount:             msg.Amount,
		PeriodDelegationId: msg.PeriodDelegationId,
	}); err != nil {
		return err
	}

	periodDel, err := k.GetPeriodDelegation(ctx, delegatorAddress, valSrcAddr, msg.PeriodDelegationId)
	if err != nil {
		return err
	}

	if err := k.ValidateDelegateMsg(ctx, &types.MsgDelegate{
		ValidatorAddress:   msg.ValidatorDstAddress,
		DelegatorAddress:   msg.DelegatorAddress,
		Amount:             msg.Amount,
		PeriodDelegationId: msg.PeriodDelegationId,
		PeriodType:         periodDel.PeriodType,
	}); err != nil {
		return err
	}

	return nil
}
