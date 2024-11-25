package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// Parameter keys
var (
	ParamStoreKeyUbi                 = []byte("ubi")
	ParamStoreKeyWithdrawAddrEnabled = []byte("withdrawaddrenabled")
	ParamStoreKeyMaxUbi              = []byte("maxubi")
)

// Deprecated: ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Deprecated: ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyUbi, &p.Ubi, validateUbi),
		paramtypes.NewParamSetPair(ParamStoreKeyWithdrawAddrEnabled, &p.WithdrawAddrEnabled, validateWithdrawAddrEnabled),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxUbi, &p.MaxUbi, validateMaxUbi),
	}
}
