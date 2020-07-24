package client

import (
	govclient "github.com/netcloth/netcloth-chain/app/v0/gov/client"
	"github.com/netcloth/netcloth-chain/app/v0/params/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/params/client/rest"
)

// ProposalHandler - param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
