package sunchain

import (
	"github.com/trinhtan/cosmos-hackathon/x/sunchain/keeper"
	"github.com/trinhtan/cosmos-hackathon/x/sunchain/types"
)

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier

	NewProduct                  = types.NewProduct
	NewMsgSetProduct            = types.NewMsgSetProduct
	NewMsgSetProductTitle       = types.NewMsgSetProductTitle
	NewMsgSetProductDescription = types.NewMsgSetProductDescription
	NewMsgDeleteProduct         = types.NewMsgDeleteProduct

	NewSell       = types.NewSell
	NewMsgSetSell = types.NewMsgSetSell

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper

	MsgBuySun          = types.MsgBuySun
	MsgSetSourceChannel = types.MsgSetSourceChannel
)
