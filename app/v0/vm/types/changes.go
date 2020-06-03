package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	// changes to the account
	createObjectChange struct {
		account *sdk.AccAddress
	}

	resetObjectChange struct {
		prev *stateObject
	}

	suicideChange struct {
		account     *sdk.AccAddress
		prev        bool // whether account had already suicided
		prevBalance sdk.Int
	}

	// changes to individual accounts
	balanceChange struct {
		account *sdk.AccAddress
		prev    sdk.Int
	}

	nonceChange struct {
		account *sdk.AccAddress
		prev    uint64
	}

	storageChange struct {
		account        *sdk.AccAddress
		key, prevValue sdk.Hash
	}

	codeChange struct {
		account            *sdk.AccAddress
		prevCode, prevHash []byte
	}

	// changes to other state values
	refundChange struct {
		prev uint64
	}

	addLogChange struct {
		txhash sdk.Hash
	}

	addPreimageChange struct {
		hash sdk.Hash
	}

	touchChange struct {
		account *sdk.AccAddress
	}
)

// createObjectChange
func (ch createObjectChange) revert(s *CommitStateDB) {
	delete(s.stateObjects, (ch.account).String())
	delete(s.stateObjectsDirty, (ch.account).String())
}

func (ch createObjectChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// resetObjectChange
func (ch resetObjectChange) revert(s *CommitStateDB) {
	s.setStateObject(ch.prev)
}

func (ch resetObjectChange) dirtied() *sdk.AccAddress {
	return nil
}

// suicideChange
func (ch suicideChange) revert(s *CommitStateDB) {
	so := s.getStateObject(*ch.account)
	if so != nil {
		so.suicided = ch.prev
		so.setBalance(ch.prevBalance)
	}
}

func (ch suicideChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// touchChange
func (ch touchChange) revert(s *CommitStateDB) {

}

func (ch touchChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// balanceChange
func (ch balanceChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setBalance(ch.prev)
}

func (ch balanceChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// nonceChange
func (ch nonceChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch nonceChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// codeChange
func (ch codeChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setCode(sdk.BytesToHash(ch.prevHash), ch.prevCode)
}

func (ch codeChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// storageChange
func (ch storageChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevValue)
}

func (ch storageChange) dirtied() *sdk.AccAddress {
	return ch.account
}

// refundChange
func (ch refundChange) revert(s *CommitStateDB) {
	s.refund = ch.prev
}

func (ch refundChange) dirtied() *sdk.AccAddress {
	return nil
}

// addLogChange
func (ch addLogChange) revert(s *CommitStateDB) {
	logs := s.logs[ch.txhash]
	if len(logs) == 1 {
		delete(s.logs, ch.txhash)
	} else {
		s.logs[ch.txhash] = logs[:len(logs)-1]
	}

	s.updateLogIndex(true)
}

func (ch addLogChange) dirtied() *sdk.AccAddress {
	return nil
}

// addPreimageChange
func (ch addPreimageChange) revert(s *CommitStateDB) {
	delete(s.preimages, ch.hash)
}

func (ch addPreimageChange) dirtied() *sdk.AccAddress {
	return nil
}
