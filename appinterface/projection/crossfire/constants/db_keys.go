package constants

import "fmt"

const VOTED_PROPOSAL_ID string = "voted_proposal_id"
const NETWORK_UPGRADE string = "network_upgrade"
const DB_KEY_SEPARATOR string = ":"
const TASK_PHASE_2_PROPOSAL_VOTE_COLUMN_NAME = "task_phase_2_proposal_vote"

const PHASE1_BLOCK_COUNT ChainStatsKey = "phase1_block_count"
const PHASE2_BLOCK_COUNT ChainStatsKey = "phase2_block_count"
const PHASE3_BLOCK_COUNT ChainStatsKey = "phase3_block_count"
type ChainStatsKey = string

const PHASE1_COMMIT CommitPhaseKey = "phase1_commit_count"
const PHASE2_COMMIT CommitPhaseKey = "phase2_commit_count"
const PHASE3_COMMIT CommitPhaseKey = "phase3_commit_count"
type CommitPhaseKey = string

func ValidatorCommitmentKey(operatorAddress string, phase CommitPhaseKey) ChainStatsKey {
	return fmt.Sprintf("%s%s%s", operatorAddress, DB_KEY_SEPARATOR, phase)
}