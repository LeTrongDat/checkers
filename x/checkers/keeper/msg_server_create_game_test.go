package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/LeTrongDat/checkers/testutil/keeper"
	"github.com/LeTrongDat/checkers/x/checkers"
	"github.com/LeTrongDat/checkers/x/checkers/keeper"
	"github.com/LeTrongDat/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = "cosmos17tsnrh8xvlk9epawu4536fun397mj9j020vkvh"
	bob   = "cosmos1669z3qdseth93n83whlx72qrnjdwdqn4gcz8xx"
	carol = "cosmos1xr0ug9d8r5yg0nay0gtsxek06yj0zmt6sldatq"
)

func setupMsgServerCreateGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), *k, sdk.WrapSDKContext(ctx)
}

func TestCreateGame(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	createGameResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     bob,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "1",
	}, *createGameResponse)
}

func TestCreateGameHasSaved(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     bob,
	})
	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, systemInfo, types.SystemInfo{
		NextId: 2,
	})

	game1, found1 := keeper.GetStoredGame(ctx, "1")
	require.True(t, found1)
	require.EqualValues(t, game1, types.StoredGame{
		Index: "1",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: alice,
		Red:   bob,
	})
}

func TestCreateGameGetAll(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     bob,
	})
	ctx := sdk.UnwrapSDKContext(context)
	games := keeper.GetAllStoredGame(ctx)
	require.Len(t, games, 1)
	require.EqualValues(t, games[0], types.StoredGame{
		Index: "1",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: alice,
		Red:   bob,
	})
}

func TestCreateGameRedAddressBad(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     "notanaddress",
	})
	require.Nil(t, createResponse)
	require.EqualError(t, err, "red address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1")
}

func TestCreateGameEmptyRedAddress(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     "",
	})
	require.Nil(t, createResponse)
	require.EqualError(t, err, "red address is invalid: : empty address string is not allowed")
}

func TestCreate3Games(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   alice,
		Red:     bob,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   bob,
		Red:     alice,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     alice,
	})
	ctx := sdk.UnwrapSDKContext(context)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, systemInfo, types.SystemInfo{
		NextId: 4,
	})

	games := keeper.GetAllStoredGame(ctx)
	require.Len(t, games, 3)

	require.EqualValues(t, games[0], types.StoredGame{
		Index: "1",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Black: alice,
		Red:   bob,
		Turn:  "b",
	})

	require.EqualValues(t, games[1], types.StoredGame{
		Index: "2",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Black: bob,
		Red:   alice,
		Turn:  "b",
	})

	require.EqualValues(t, games[2], types.StoredGame{
		Index: "3",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Black: bob,
		Red:   alice,
		Turn:  "b",
	})
}

func TestCreateGameFarFuture(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	systemInfo, found := keeper.GetSystemInfo(ctx)
	systemInfo.NextId = 1024
	keeper.SetSystemInfo(ctx, systemInfo)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     bob,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "1024",
	}, *createResponse)

	systemInfo, found = keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, systemInfo, types.SystemInfo{
		NextId: 1025,
	})

	game, found := keeper.GetStoredGame(ctx, "1024")
	require.True(t, found)
	require.EqualValues(t, game, types.StoredGame{
		Index: "1024",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: alice,
		Red:   bob,
	})

}
