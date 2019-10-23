package types

import (
	"fmt"
	"github.com/NetCloth/netcloth-chain/codec"
	sdk "github.com/NetCloth/netcloth-chain/types"
	"time"
)

type UnstakingEntry struct {
	EndTime             time.Time   `json:"end_time" yaml:"end_time"`
	Amount              sdk.Coin    `json:"amount" yaml:"amount"`
}

func NewUnstakingEntry(creationHeight int64, endTime time.Time, amount sdk.Coin) UnstakingEntry {
	return UnstakingEntry {
		EndTime: endTime,
		Amount:  amount,
	}
}

func (e UnstakingEntry) IsMature(currentTime time.Time) bool {
	return !e.EndTime.After(currentTime)
}

type Unstaking struct {
	AccountAddress      sdk.AccAddress      `json:"account_address" yaml:"account_address"`
	Entries             []UnstakingEntry    `json:"entries" yaml:"entries"`
}

func NewUnstaking(accAddress sdk.AccAddress, creationHeight int64, completionTime time.Time, balance sdk.Coin) Unstaking  {
	entry := NewUnstakingEntry(creationHeight, completionTime, balance)
	return Unstaking{
		AccountAddress: accAddress,
		Entries:        []UnstakingEntry{entry},
	}
}

func (s *Unstaking) AddEntry(creationHeight int64, minTime time.Time, amount sdk.Coin) {
	entry := NewUnstakingEntry(creationHeight, minTime, amount)
	s.Entries = append(s.Entries, entry)
}

func (s *Unstaking) RemoveEntry(i int64) {
	s.Entries = append(s.Entries[:i], s.Entries[i+1:]...)
}

func MustMarshalUnstaking(cdc *codec.Codec, unstaking Unstaking) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(unstaking)
}

func MustUnmarshalUnstaking(cdc *codec.Codec, value []byte) Unstaking {
	unstaking, err := UnmarshalUnstaking(cdc, value)
	if err != nil {
		panic(err)
	}
	return unstaking
}

func UnmarshalUnstaking(cdc *codec.Codec, value []byte) (unstaking Unstaking, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &unstaking)
	return unstaking, err
}

func (s Unstaking) String() string {
	out := fmt.Sprintf(`Unstaking:
  AccountAddress: %s
  Entries:`, s.AccountAddress)

	for i, entry := range s.Entries {
		out += fmt.Sprintf(`    Unstaking %d:
      EndTime   : %v
      Amount    : %s`, i, entry.EndTime.String(), entry.Amount)
	}

	return out
}

type UnStakingTODO struct {
	AccountAddress sdk.AccAddress `json:"account_address" yaml:"account_address"`
	Amount sdk.Coin `json:"amount" yaml:"amount"`
	EndTime time.Time `json:"end_time" yaml:"end_time"`
}

func NewUnStakingTODO(accountAddress sdk.AccAddress, amount sdk.Coin, endTime time.Time) UnStakingTODO {
	return UnStakingTODO {
		AccountAddress: accountAddress,
		Amount: amount,
		EndTime:endTime,
	}
}

