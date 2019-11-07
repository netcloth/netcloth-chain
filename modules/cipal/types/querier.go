package types

const (
	QueryIPAL = "query"
)

type QueryCIPALParams struct {
	AccAddr string
}

func NewQueryCIPALParams(AccAddr string) QueryCIPALParams {
	return QueryCIPALParams{
		AccAddr: AccAddr,
	}
}
