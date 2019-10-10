package types

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "ipal"

	// StoreKey is the store key string for ipal
	StoreKey = ModuleName

	// RouterKey is the message route for ipal
	RouterKey = ModuleName

	// QuerierRoute is the querier route for ipal
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

var (
	IPALObjectKey = []byte{0x11} // prefix for each key to an ipal object
)

func GetIPALObjectKey(addr string) []byte {
	return append(IPALObjectKey, []byte(addr)...)
}