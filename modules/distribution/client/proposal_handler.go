package client

import (
	"github.com/netcloth/netcloth-chain/modules/distribution/client/cli"
	"github.com/netcloth/netcloth-chain/modules/distribution/client/rest"
	govclient "github.com/netcloth/netcloth-chain/modules/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
