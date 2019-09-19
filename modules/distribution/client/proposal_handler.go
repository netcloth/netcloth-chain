package client

import (
	"github.com/NetCloth/netcloth-chain/modules/distribution/client/cli"
	"github.com/NetCloth/netcloth-chain/modules/distribution/client/rest"
	govclient "github.com/NetCloth/netcloth-chain/modules/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
