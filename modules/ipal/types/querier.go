package types

const (
	QueryIPAL = "ipal"
)

type QueryIPALParams struct {
	AccAddr string
}

func NewQueryIPALParams(AccAddr string) QueryIPALParams {
	return QueryIPALParams{
		AccAddr: AccAddr,
	}
}
