package types

const (
	// ModuleName is the name of the vm module
	ModuleName = "vm"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the vm module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the vm module
	RouterKey = ModuleName
)

var (
	ContractCodeKey = []byte{0x10}
)

func GetContractCodeKey(contracAddr, codeHash []byte) []byte {
	v := append(ContractCodeKey, contracAddr...)
	return append(v, codeHash[:30]...)
}
