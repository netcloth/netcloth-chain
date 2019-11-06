package types

const (
	QueryIPAL = "cipal"
)

type QueryCIPALParams struct {
	AccAddr string
}

func NewQueryCIPALParams(AccAddr string) QueryCIPALParams {
	return QueryCIPALParams{
		AccAddr: AccAddr,
	}
}
