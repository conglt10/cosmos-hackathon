package sunchain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"

	"github.com/trinhtan/cosmos-hackathon/x/sunchain/types"
)


// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgBuySun:
			return handleBuySun(ctx, msg, keeper)
		case MsgSetSourceChannel:
			return handleSetSourceChannel(ctx, msg, keeper)
		case channeltypes.MsgPacket:
			var responsePacket oracle.OracleResponsePacketData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &responsePacket); err == nil {
				return handleOracleRespondPacketData(ctx, responsePacket, keeper)
			}
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal oracle packet data")
		case MsgSetProduct:
			return handleMsgSetProduct(ctx, keeper, msg)
		case MsgSetProductTitle:
			return handleMsgSetProductTitle(ctx, keeper, msg)
		case MsgSetProductDescription:
			return handleMsgSetProductDescription(ctx, keeper, msg)
		case MsgDeleteProduct:
			return handleMsgDeleteProduct(ctx, keeper, msg)
		case MsgSetSell:
			return handleMsgSetSell(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}

func handleBuySun(ctx sdk.Context, msg MsgBuyGold, keeper Keeper) (*sdk.Result, error) {
	orderID, err := keeper.AddOrder(ctx, msg.Buyer, msg.Amount)
	if err != nil {
		return nil, err
	}
	// TODO: Set all bandchain parameter here
	bandChainID := "bandchain"
	port := "sunchain"
	oracleScriptID := oracle.OracleScriptID(3)
	calldata := make([]byte, 8)
	binary.LittleEndian.PutUint64(calldata, 1000000)
	askCount := int64(1)
	minCount := int64(1)

	channelID, err := keeper.GetChannel(ctx, bandChainID, port)

	if err != nil {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"not found channel to bandchain",
		)
	}
	sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, port, channelID)
	if !found {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"unknown channel %s port sunchain",
			channelID,
		)
	}
	destinationPort := sourceChannelEnd.Counterparty.PortID
	destinationChannel := sourceChannelEnd.Counterparty.ChannelID
	sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(
		ctx, port, channelID,
	)
	if !found {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"unknown sequence number for channel %s port oracle",
			channelID,
		)
	}
	packet := oracle.NewOracleRequestPacketData(
		fmt.Sprintf("Order:%d", orderID), oracleScriptID, hex.EncodeToString(calldata),
		askCount, minCount,
	)
	err = keeper.ChannelKeeper.SendPacket(ctx, channel.NewPacket(packet.GetBytes(),
		sequence, port, channelID, destinationPort, destinationChannel,
		1000000000, // Arbitrarily high timeout for now
	))
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleSetSourceChannel(ctx sdk.Context, msg MsgSetSourceChannel, keeper Keeper) (*sdk.Result, error) {
	keeper.SetChannel(ctx, msg.ChainName, msg.SourcePort, msg.SourceChannel)
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleOracleRespondPacketData(ctx sdk.Context, packet oracle.OracleResponsePacketData, keeper Keeper) (*sdk.Result, error) {
	clientID := strings.Split(packet.ClientID, ":")
	if len(clientID) != 2 {
		return nil, sdkerrors.Wrapf(types.ErrUnknownClientID, "unknown client id %s", packet.ClientID)
	}
	orderID, err := strconv.ParseUint(clientID[1], 10, 64)
	if err != nil {
		return nil, err
	}
	rawResult, err := hex.DecodeString(packet.Result)
	if err != nil {
		return nil, err
	}
	result, err := types.DecodeResult(rawResult)
	if err != nil {
		return nil, err
	}

	// Assume multiplier should be 1000000
	order, err := keeper.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	// TODO: Calculate collateral percentage
	goldAmount := order.Amount[0].Amount.Int64() / int64(result.Px)
	if goldAmount == 0 {
		escrowAddress := types.GetEscrowAddress()
		err = keeper.BankKeeper.SendCoins(ctx, escrowAddress, order.Owner, order.Amount)
		if err != nil {
			return nil, err
		}
		order.Status = types.Completed
		keeper.SetOrder(ctx, orderID, order)
	} else {
		goldToken := sdk.NewCoin("gold", sdk.NewInt(goldAmount))
		keeper.BankKeeper.AddCoins(ctx, order.Owner, sdk.NewCoins(goldToken))
		order.Gold = goldToken
		order.Status = types.Active
		keeper.SetOrder(ctx, orderID, order)
	}
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}


// handleMsgProduct handles a message to set product
func handleMsgSetProduct(ctx sdk.Context, keeper Keeper, msg MsgSetProduct) (*sdk.Result, error) {
	// if !msg.Owner.Equals(keeper.GetProductOwner(ctx, msg.ProductID)) { // Checks if the the msg sender is the same as the current owner
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	// }

	var product = Product{
		ProductID:   msg.ProductID,
		Title:       msg.Title,
		Description: msg.Description,
		Owner:       msg.Owner,
	}

	keeper.SetProduct(ctx, msg.ProductID, product)
	return &sdk.Result{}, nil // return
}

// handleMsgSetProductTitle handles a message to set product title
func handleMsgSetProductTitle(ctx sdk.Context, keeper Keeper, msg MsgSetProductTitle) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetProductOwner(ctx, msg.ProductID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	keeper.SetProductTitle(ctx, msg.ProductID, msg.Title)
	return &sdk.Result{}, nil // return
}

// handleMsgSetProductDescription handles a message to set product description
func handleMsgSetProductDescription(ctx sdk.Context, keeper Keeper, msg MsgSetProductDescription) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetProductOwner(ctx, msg.ProductID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	keeper.SetProductDescription(ctx, msg.ProductID, msg.Description)
	return &sdk.Result{}, nil // return
}

// Handle a message to delete product
func handleMsgDeleteProduct(ctx sdk.Context, keeper Keeper, msg MsgDeleteProduct) (*sdk.Result, error) {
	if !keeper.IsProductPresent(ctx, msg.ProductID) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.ProductID)
	}
	if !msg.Owner.Equals(keeper.GetProductOwner(ctx, msg.ProductID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.DeleteWhois(ctx, msg.ProductID)
	return &sdk.Result{}, nil
}

// handleMsgSetSell handles a message to set sell
func handleMsgSetSell(ctx sdk.Context, keeper Keeper, msg MsgSetSell) (*sdk.Result, error) {

	var sell = Sell{
		SellID:    msg.SellID,
		ProductID: msg.ProductID,
		Seller:    msg.Seller,
		MinPrice:  msg.MinPrice,
	}

	keeper.SetSell(ctx, msg.SellID, sell)
	return &sdk.Result{}, nil // return
}
