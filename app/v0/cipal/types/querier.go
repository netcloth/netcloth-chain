package types

const (
	QueryCIPAL      = "query"
	QueryCIPALCount = "count"
	QueryCIPALs     = "batch_query"
)

type QueryCIPALParams struct {
	AccAddr string
}

type QueryCIPALsParams struct {
	AccAddrs []string `json:"acc_addrs"`
}

func NewQueryCIPALParams(accAddr string) QueryCIPALParams {
	return QueryCIPALParams{
		AccAddr: accAddr,
	}
}
