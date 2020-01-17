package types

// Handler defines the core of the state transition function of an application.
type Handler func(ctx Context, msg Msg) (*Result, error)

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx Context, tx Tx, simulate bool) (newCtx Context, err error)
type FeePreprocessHandler func(ctx Context, tx Tx) Error
type FeeRefundHandler func(ctx Context, tx Tx, result Result) (Coin, error)
