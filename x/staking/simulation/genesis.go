package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Simulation parameter constants
const (
	unbondingTime     = "unbonding_time"
	maxValidators     = "max_validators"
	historicalEntries = "historical_entries"
)

// genUnbondingTime returns randomized UnbondingTime
func genUnbondingTime(r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// genMaxValidators returns randomized MaxValidators
func genMaxValidators(r *rand.Rand) (maxValidators uint32) {
	return uint32(r.Intn(250) + 1)
}

// getHistEntries returns randomized HistoricalEntries between 0-100.
func getHistEntries(r *rand.Rand) uint32 {
	return uint32(r.Intn(int(types.DefaultHistoricalEntries + 1)))
}

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var (
		unbondTime        time.Duration
		maxVals           uint32
		histEntries       uint32
		minCommissionRate sdkmath.LegacyDec
	)

	simState.AppParams.GetOrGenerate(unbondingTime, &unbondTime, simState.Rand, func(r *rand.Rand) { unbondTime = genUnbondingTime(r) })

	simState.AppParams.GetOrGenerate(maxValidators, &maxVals, simState.Rand, func(r *rand.Rand) { maxVals = genMaxValidators(r) })

	simState.AppParams.GetOrGenerate(historicalEntries, &histEntries, simState.Rand, func(r *rand.Rand) { histEntries = getHistEntries(r) })

	// NOTE: the slashing module need to be defined after the staking module on the
	// NewSimulationManager constructor for this to work
	simState.UnbondTime = unbondTime
	params := types.NewParams(simState.UnbondTime, maxVals, 7, histEntries, simState.BondDenom, minCommissionRate, types.DefaultMinDelegation, types.DefaultFlexiblePeriodType, types.DefaultPeriods, types.DefaultLockedTokenType, types.DefaultTokenTypes, types.DefaultSingularityHeight)

	// validators & delegations
	var (
		validators        []types.Validator
		delegations       []types.Delegation
		periodDelegations []types.PeriodDelegation
	)

	valAddrs := make([]sdk.ValAddress, simState.NumBonded)

	now := time.Now()

	for i := 0; i < int(simState.NumBonded); i++ {
		valAddr := sdk.ValAddress(simState.Accounts[i].Address)
		valAddrs[i] = valAddr

		maxCommission := sdkmath.LegacyNewDecWithPrec(int64(simulation.RandIntBetween(simState.Rand, 1, 100)), 2)
		commission := types.NewCommission(
			simulation.RandomDecAmount(simState.Rand, maxCommission),
			maxCommission,
			simulation.RandomDecAmount(simState.Rand, maxCommission),
		)

		validator, err := types.NewValidator(valAddr.String(), simState.Accounts[i].ConsKey.PubKey(), types.Description{}, 0)
		if err != nil {
			panic(err)
		}
		validator.Tokens = simState.InitialStake
		validator.DelegatorShares = sdkmath.LegacyNewDecFromInt(simState.InitialStake)
		validator.Commission = commission

		delegation := types.NewDelegation(
			simState.Accounts[i].Address.String(), valAddr.String(), sdkmath.LegacyNewDecFromInt(simState.InitialStake), sdkmath.LegacyNewDecFromInt(simState.InitialStake),
		)
		periodDelegation := types.NewPeriodDelegation(
			simState.Accounts[i].Address.String(),
			valAddr.String(),
			"1",
			sdkmath.LegacyNewDecFromInt(simState.InitialStake),
			sdkmath.LegacyNewDecFromInt(simState.InitialStake),
			1,
			now.Add(time.Hour*24*30),
		)

		validators = append(validators, validator)
		delegations = append(delegations, delegation)
		periodDelegations = append(periodDelegations, periodDelegation)
	}

	stakingGenesis := types.NewGenesisState(params, validators, delegations, periodDelegations)

	bz, err := json.MarshalIndent(&stakingGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated staking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(stakingGenesis)
}
