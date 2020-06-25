package nameservice

import (
	"github.com/althea-net/peggy/module/x/nameservice/keeper"
	"github.com/althea-net/peggy/module/x/nameservice/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.ModuleName
	QuerierRoute      = types.QuerierRoute
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewMsgBuyName       = types.NewMsgBuyName
	NewMsgSetName       = types.NewMsgSetName
	NewMsgDeleteName    = types.NewMsgDeleteName
	NewMsgSetEthAddress = types.NewMsgSetEthAddress
	NewWhois            = types.NewWhois
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
)

type (
	Keeper           = keeper.Keeper
	MsgSetName       = types.MsgSetName
	MsgBuyName       = types.MsgBuyName
	MsgDeleteName    = types.MsgDeleteName
	MsgSetEthAddress = types.MsgSetEthAddress
	QueryResResolve  = types.QueryResResolve
	QueryResNames    = types.QueryResNames
	Whois            = types.Whois
)
