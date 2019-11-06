package types

const (
	ModuleName   = "cipal"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	CIPALObjectKey = []byte{0x11}
)

func GetCIPALObjectKey(addr string) []byte {
	return append(CIPALObjectKey, []byte(addr)...)
}
