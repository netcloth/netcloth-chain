package types

const (
	ModuleName = "ipal"
	StoreKey = ModuleName
	RouterKey = ModuleName
	QuerierRoute = ModuleName
)

var (
	IPALObjectKey = []byte{0x11}
)

func GetIPALObjectKey(addr string) []byte {
	return append(IPALObjectKey, []byte(addr)...)
}
