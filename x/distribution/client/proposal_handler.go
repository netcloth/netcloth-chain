package client

import (
	"github.com/NetCloth/netcloth-chain/x/distribution/client/cli"
	"github.com/NetCloth/netcloth-chain/x/distribution/client/rest"
	govclient "github.com/NetCloth/netcloth-chain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
