package context

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	tmlite "github.com/tendermint/tendermint/lite"
	tmliteproxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var (
	verifierDir = ".lite_verifier"

	// DefaultVerifierCacheSize defines the default Tendermint cache size.
	DefaultVerifierCacheSize = 10
)

func CreateVerifier(ctx CLIContext, cacheSize int) (tmlite.Verifier, error) {
	if ctx.TrustNode {
		return nil, nil
	}

	switch {
	case ctx.ChainID == "":
		return nil, errors.New("must provide a valid chain ID to create verifier")

	case ctx.HomeDir == "":
		return nil, errors.New("must provide a valid home directory to create verifier")

	case ctx.Client == nil && ctx.NodeURI == "":
		return nil, errors.New("must provide a valid RPC client or RPC URI to create verifier")
	}

	var err error

	// create an RPC client based off of the RPC URI if no RPC client exists
	client := ctx.Client
	if client == nil {
		client, err = rpcclient.NewHTTP(ctx.NodeURI, "/websocket")
		if err != nil {
			return nil, err
		}
	}

	return tmliteproxy.NewVerifier(
		ctx.ChainID, filepath.Join(ctx.HomeDir, ctx.ChainID, verifierDir),
		client, log.NewNopLogger(), cacheSize,
	)
}
