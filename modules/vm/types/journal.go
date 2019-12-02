package types

import sdk "github.com/netcloth/netcloth-chain/types"

// journalEntry is a modification entry in the state change journal that can be reverted on demand
type journalEntry interface {
	// revert undoes the changes introduced by this journal entry
	revert(*CommitStateDB)

	// dirtied returns the address modified by this journal entry
	dirtied() *sdk.AccAddress
}

// journal contains the list of state modifications applied since the last state
// commit. These are tracked to be able to be reverted in case of an execution
// exception or revert request
type journal struct {
	entries []journalEntry // current changes tracked by the journal
	dirties map[string]int // dirty accounts and the number of changes
	//TODO: map key --> sdk.AccAddress
}

func newJournal() *journal {
	return &journal{
		dirties: make(map[string]int),
	}
}

// append inserts a new modification entry to the end of the change journal
func (j *journal) append(entry journalEntry) {
	j.entries = append(j.entries, entry)
	if addr := entry.dirtied(); addr != nil {
		j.dirties[(*addr).String()]++
	}
}

func (j *journal) revert(statedb *CommitStateDB, snapshot int) {
	for i := len(j.entries) - 1; i >= snapshot; i-- {
		// undo the changes made by the operation
		j.entries[i].revert(statedb)

		// drop any dirty tracking induced by the change
		if addr := j.entries[i].dirtied(); addr != nil {
			if j.dirties[(*addr).String()]--; j.dirties[(*addr).String()] == 0 {
				delete(j.dirties, (*addr).String())
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

func (j *journal) dirty(addr sdk.AccAddress) {
	j.dirties[addr.String()]++
}

func (j *journal) length() int {
	return len(j.entries)
}
