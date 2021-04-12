package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Question: what here can be epoched?
	k.slash(ctx)
	k.tallyAttestations(ctx)
	k.timeoutTxs(ctx)
	k.createEthSignerSet(ctx)
}

// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func (k Keeper) tallyAttestations(ctx sdk.Context) {
	// We check the attestations that haven't been observed, i.e nonce is exactly 1 higher than the last attestation
	nonce := uint64(k.GetLastObservedEventNonce(ctx)) + 1

	k.IterateAttestationByNonce(ctx, nonce, func(attestation types.Attestation) bool {
		// try unobserved attestations
		// TODO: rename. "Try" is too ambiguous
		k.TryAttestation(ctx, attestation)
		return false
	})
}

// timeoutTxs deletes the batch and logic call transactions that have passed
// their expiration height on Ethereum.
func (k Keeper) timeoutTxs(ctx sdk.Context) {
	info, found := k.GetEthereumInfo(ctx)
	if !found {
		panic("ethereum observed info not found")
	}

	// TODO: start iteration in desc order from height = info.Height
	// TODO: can we iterate once for over a height range [0, info.Height] instead of
	// once for every tx type
	k.IterateBatchTxs(ctx, func(tx types.BatchTx) bool {
		if tx.Timeout < info.Height {
			k.CancelBatchTx(ctx, tx.TokenContract, tx.Nonce)
		}

		return false
	})

	k.IterateLogicCallTxs(ctx, func(tx types.LogicCallTx) bool {
		if tx.Timeout < info.Height {
			k.CancelLogicCallTx(ctx, tx.InvalidationId, tx.InvalidationNonce)
		}

		return false
	})
}

// Auto ValsetRequest Creation.
// 1. If there are no valset requests, create a new one. TODO: why? is it necessary?
// 2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
//    This will make sure the unbonding validator has to provide an attestation to a new Valset
//      that excludes him before he completely Unbonds.  Otherwise he will be slashed
// 3. If power change between validators of CurrentValset and latest valset request is > 5%
func (k Keeper) createEthSignerSet(ctx sdk.Context) {

	latestValset := k.GetLatestValset(ctx)
	lastUnbondingHeight := k.GetLastUnbondingBlockHeight(ctx)

	if (latestValset == nil) || (lastUnbondingHeight == uint64(ctx.BlockHeight())) || (types.BridgeValidators(k.GetCurrentValset(ctx).Members).PowerDiff(latestValset.Members) > 0.05) {
		k.SetValsetRequest(ctx)
	}
}