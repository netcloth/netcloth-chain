package client

import (
	"github.com/NetCloth/netcloth-chain/x/params/client/cli"
	"github.com/NetCloth/netcloth-chain/x/params/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
