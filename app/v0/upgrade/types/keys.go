package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/netcloth/netcloth-chain/app/protocol"
)

const (
	ModuleName   = protocol.UpgradeModuleName
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	proposalIDKey     = "p/%s"         // p/<proposalId>
	successVersionKey = "success/%s"   // success/<protocolVersion>
	failedVersionKey  = "failed/%s/%s" // failed/<protocolVersion>/<proposalId>
	signalKey         = "s/%s/%s"      // s/<protocolVersion>/<switchVoterAddress>
	signalPrefixKey   = "s/%s"
)

// GetProposalIDKey gets proposal ID store key
func GetProposalIDKey(proposalID uint64) []byte {
	return []byte(fmt.Sprintf(proposalIDKey, UintToHexString(proposalID)))
}

// GetSuccessVersionKey gets successful version store key
func GetSuccessVersionKey(versionID uint64) []byte {
	return []byte(fmt.Sprintf(successVersionKey, UintToHexString(versionID)))
}

// GetFailedVersionKey gets failed version store key
func GetFailedVersionKey(versionID uint64, proposalID uint64) []byte {
	return []byte(fmt.Sprintf(failedVersionKey, UintToHexString(versionID), UintToHexString(proposalID)))
}

// GetSignalKey gets signal store key
func GetSignalKey(versionID uint64, switchVoterAddr string) []byte {
	return []byte(fmt.Sprintf(signalKey, UintToHexString(versionID), switchVoterAddr))
}

// GetSignalPrefixKey gets signal prefix store key
func GetSignalPrefixKey(versionID uint64) []byte {
	return []byte(fmt.Sprintf(signalPrefixKey, UintToHexString(versionID)))
}

// GetAddressFromSignalKey gets address from signal key
func GetAddressFromSignalKey(key []byte) string {
	return strings.Split(string(key), "/")[2]
}

// UintToHexString converts uint to hex string
func UintToHexString(i uint64) string {
	hex := strconv.FormatUint(i, 16)
	var stringBuild bytes.Buffer
	for i := 0; i < 16-len(hex); i++ {
		stringBuild.Write([]byte("0"))
	}
	stringBuild.Write([]byte(hex))
	return stringBuild.String()
}
