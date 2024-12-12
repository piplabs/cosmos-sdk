package v3_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	v3 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v3"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestMigrateJSON(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithCodec(encodingConfig.Codec)

	oldState := types.DefaultGenesisState()

	newState, err := v3.MigrateJSON(*oldState)
	require.NoError(t, err)

	bz, err := clientCtx.Codec.MarshalJSON(&newState)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about new param MinCommissionRate.
	expected := `{
	"delegations": [],
	"exported": false,
	"last_total_power": "0",
	"last_validator_powers": [],
	"params": {
		"bond_denom": "stake",
		"flexible_period_type": 0,
		"historical_entries": 10000,
		"locked_token_type": 0,
		"max_entries": 7,
		"max_validators": 100,
		"min_commission_rate": "0.000000000000000000",
		"min_delegation": "1",
		"periods": [
			{
				"duration": "0s",
				"period_type": 0,
				"rewards_multiplier": "1.000000000000000000"
			},
			{
				"duration": "7776000s",
				"period_type": 1,
				"rewards_multiplier": "1.051000000000000000"
			},
			{
				"duration": "31104000s",
				"period_type": 2,
				"rewards_multiplier": "1.160000000000000000"
			},
			{
				"duration": "46656000s",
				"period_type": 3,
				"rewards_multiplier": "1.340000000000000000"
			}
		],
		"singularity_height": "403200",
		"token_types": [
			{
				"rewards_multiplier": "0.500000000000000000",
				"token_type": 0
			},
			{
				"rewards_multiplier": "1.000000000000000000",
				"token_type": 1
			}
		],
		"unbonding_time": "1814400s"
	},
	"period_delegations": [],
	"redelegations": [],
	"unbonding_delegations": [],
	"validators": []
}`

	require.Equal(t, expected, string(indentedBz))
}
