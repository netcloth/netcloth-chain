package token

import (
	"fmt"
	nch "github.com/NetCloth/netcloth-chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"strings"
)

var (
	MinimumMonikerSize = 3
	MaximumMonikerSize = 8


	IsAlphaNumeric     = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString   // only accepts alphanumeric characters
	IsAlphaNumericDash = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString // only accepts alphanumeric characters, _ and -
	IsBeginWithAlpha   = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
)

// NewHandler returns a handler
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle handleMsgIssue
func handleMsgIssue(ctx sdk.Context, k Keeper, msg MsgIssue) sdk.Result {
	// check issue amount
	if !msg.Amount.IsValid() {
		return sdk.ErrInsufficientCoins("invalid coins").Result()
	}

	// check moniker
	if err := ValidateMoniker(msg.Amount.Denom); err != nil {
		return err.Result()
	}

	ctx.Logger().Debug(
		fmt.Sprintf(
			"issue coins, from: %s,  to: %s, amount: %s ",
			msg.Banker.String(), msg.Amount.String(), msg.Address.String(),
			),
		)

	newCoins := sdk.NewCoins(msg.Amount)
	_, err := k.coinKeeper.AddCoins(ctx, msg.Address, newCoins)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func ValidateMoniker(moniker string) sdk.Error {
	// check the moniker size
	if len(moniker) < MinimumMonikerSize || len(moniker) > MaximumMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be between [%d,%d]", MinimumMonikerSize, MaximumMonikerSize))
	}

	// check the moniker format
	if !IsBeginWithAlpha(moniker) || !IsAlphaNumeric(moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must begin with a letter followed by alphanumeric characters"))
	}

	// check if the moniker contains the native token name
	if strings.Contains(strings.ToLower(moniker), nch.NativeTokenName) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must not contain the native token name"))
	}

	// check moniker not exists
	return nil
}