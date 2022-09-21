package keeper_test

import (
	"testing"

	"github.com/LeTrongDat/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = "cosmos17tsnrh8xvlk9epawu4536fun397mj9j020vkvh"
	bob   = "cosmos1669z3qdseth93n83whlx72qrnjdwdqn4gcz8xx"
)

func TestCreateGame(t *testing.T) {
	msgServer, context := setupMsgServer(t)
	createGameResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   alice,
		Red:     bob,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "",
	}, createGameResponse)
}
