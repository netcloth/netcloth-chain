package client

import (
	"github.com/netcloth/netcloth-chain/app/v0/distribution/client/cli"
	"github.com/netcloth/netcloth-chain/app/v0/distribution/client/rest"
	govclient "github.com/netcloth/netcloth-chain/app/v0/gov/client"
)

var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
