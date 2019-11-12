package client

import (
	"github.com/netcloth/netcloth-chain/modules/params/client/cli"
	"github.com/netcloth/netcloth-chain/modules/params/client/rest"
	govclient "github.com/netcloth/netcloth-chain/modules/gov/client"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
