package client

import (
	"github.com/netcloth/netcloth-chain/app/v0/distribution/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/distribution/client/rest"
	govclient "github.com/netcloth/netcloth-chain/app/v0/gov/client"
)

// ProposalHandler - param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
