package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestIPALUserRequestSig(t *testing.T) {
	userAddress := "nch1khneh5rr978lv6rz2f55aj6u83s5efrky37477"
	serviceAddress := "nch10jzpt32gwradv9mcnr6fuuj0tnx7rq0psmmtju"
	serviceType := uint64(1)

	expireStr := "2020-03-12T10:38:55"
	expiration, err := time.ParseInLocation("2006-01-02T15:04:05", expireStr, time.UTC)
	fmt.Println(expiration)
	require.Nil(t, err)

	adParam := NewADParam(userAddress, serviceAddress, serviceType, expiration)
	//fmt.Println(fmt.Sprintf("param: %x", adParam.GetSignBytes()))

	// parse private key
	rawPrivKey := sdk.FromHex("a869398c1d25422b6bda112c91552a1de9464411767765d9a54fea44899a7f3d")
	var privKeyBytes [32]byte
	copy(privKeyBytes[:], rawPrivKey[:])
	privKey := secp256k1.PrivKeySecp256k1(privKeyBytes)

	// sign
	sig, _ := privKey.Sign(adParam.GetSignBytes())
	pub := privKey.PubKey()
	stdSig := auth.StdSignature{
		PubKey:    pub,
		Signature: sig,
	}
	//fmt.Println(fmt.Sprintf("signature: %x", stdSig.Signature))

	// verify signature
	sigVerifyPass := stdSig.VerifyBytes(adParam.GetSignBytes(), stdSig.Signature)
	require.True(t, sigVerifyPass)
}
