package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/crypto"
	yaml "gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/modules/auth/exported"
	sdk "github.com/netcloth/netcloth-chain/types"
)

//-----------------------------------------------------------------------------
// BaseAccount

var _ exported.Account = (*BaseAccount)(nil)
var _ exported.GenesisAccount = (*BaseAccount)(nil)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`

	// contract code hash
	CodeHash []byte
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address sdk.AccAddress, coins sdk.Coins,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coins:         coins,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

// ProtoBaseAccount - a prototype function for BaseAccount
func ProtoBaseAccount() exported.Account {
	return &BaseAccount{}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr sdk.AccAddress) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// GetCodeHash - Implements sdk.Account.
func (acc *BaseAccount) GetCodeHash() []byte {
	return acc.CodeHash
}

// SetCodeHash - Implements sdk.Account.
func (acc *BaseAccount) SetCodeHash(codeHash []byte) {
	acc.CodeHash = codeHash
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return acc.GetCoins()
}

// Balance returns the balance of an account
func (acc *BaseAccount) Balance() sdk.Int {
	return acc.GetCoins().AmountOf(sdk.NativeTokenName)
}

// SetBalance sets an account's balance of native token
func (acc *BaseAccount) SetBalance(amt sdk.Int) {
	coins := acc.GetCoins()
	diff := amt.Sub(coins.AmountOf(sdk.NativeTokenName))
	if diff.IsZero() {
		return
	} else if diff.IsPositive() {
		// Increase coins to amount
		coins = coins.Add(sdk.NewCoin(sdk.NativeTokenName, diff))
	} else {
		// Decrease coin to amount
		coins = coins.Sub(sdk.Coins{sdk.NewCoin(sdk.NativeTokenName, diff.Neg())})
	}

	if err := acc.SetCoins(coins); err != nil {
		panic(fmt.Sprintf("Could not set coins for address %s", acc.GetAddress()))
	}
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.PubKey != nil && acc.Address != nil &&
		!bytes.Equal(acc.PubKey.Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

type baseAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	CodeHash      []byte         `json:"code_hash" yaml:"code_hash"`
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	alias := baseAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      acc.CodeHash,
	}

	if acc.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

// MarshalJSON returns the JSON representation of a BaseAccount.
func (acc BaseAccount) MarshalJSON() ([]byte, error) {
	alias := baseAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      acc.CodeHash,
	}

	if acc.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	return json.Marshal(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into a BaseAccount.
func (acc *BaseAccount) UnmarshalJSON(bz []byte) error {
	var alias baseAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	if alias.PubKey != "" {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, alias.PubKey)
		if err != nil {
			return err
		}

		acc.PubKey = pk
	}

	acc.Address = alias.Address
	acc.Coins = alias.Coins
	acc.AccountNumber = alias.AccountNumber
	acc.Sequence = alias.Sequence
	acc.CodeHash = alias.CodeHash

	return nil
}
