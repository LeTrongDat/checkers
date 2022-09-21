package types

import (
	"errors"
	"fmt"

	"github.com/LeTrongDat/checkers/x/checkers/rules"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (storedGame StoredGame) GetBlackAddress() (black sdk.AccAddress, err error) {
	black, errBlack := sdk.AccAddressFromBech32(storedGame.Black)
	return black, sdkerrors.Wrapf(errBlack, ErrInvalidBlack.Error(), storedGame.Black)
}

func (storedGame StoredGame) GetRedAddress() (red sdk.AccAddress, err error) {
	red, errRed := sdk.AccAddressFromBech32(storedGame.Red)
	return red, sdkerrors.Wrapf(errRed, ErrInvalidRed.Error(), storedGame.Red)
}

func (storedGame StoredGame) ParseGame() (game *rules.Game, err error) {
	game, errGame := rules.Parse(storedGame.Board)
	if errGame != nil {
		return nil, sdkerrors.Wrapf(errGame, ErrGameNotParseable.Error())
	}
	game.Turn = rules.StringPieces[storedGame.Turn].Player
	if game.Turn.Color == "NO_PLAYER" {
		return nil, sdkerrors.Wrapf(errors.New(fmt.Sprintf("Turn: %s", storedGame.Turn)), ErrGameNotParseable.Error())
	}
	return game, nil
}

func (storedGame StoredGame) Validate() (err error) {
	if _, err = storedGame.GetBlackAddress(); err != nil {
		return err
	}

	if _, err = storedGame.GetRedAddress(); err != nil {
		return err
	}

	_, err = storedGame.ParseGame()
	return err
}
