package game313

import (
	"github.com/Bajron/three-thirteen/playingcards"
	"testing"
)

func TestPublicState(t *testing.T) {
	var value interface{}
	state := New(3)
	value = GetPublicStateView(state)
	public, ok := value.(*PublicStateView)
	if !ok {
		t.Error("returned value is not of type PublicStateView*")
	}

	if len(public.Players) != len(state.Players) {
		t.Error("players view count should match players' count")
	}
}

func TestGame(t *testing.T) {
	session := NewGameSession(3, 0)

	if len(session.state.Players) != 3 {
		t.Error("session should have initialized state")
	}

	mErr := session.Marshal(GameCommand{session.state.StartingPlayer, DEAL})
	if mErr == nil {
		t.Error("deal should fail for wrong player")
	}

	mErr = session.Marshal(GameCommand{session.state.DealingPlayer(), DEAL})
	if mErr != nil {
		t.Error("deal should not fail")
	}

	current := session.state.CurrentPlayer
	cmd := NewTakeCommand(current, TAKE_FROM_PILE)

	preState := GetPublicStateView(session.state)
	card, err := session.Dispatch(cmd)
	postState := GetPublicStateView(session.state)

	if card != preState.PileTop || err != nil {
		t.Error("session should forward moves to state")
	}

	if postState.PileTop != playingcards.NIL_CARD {
		t.Error("state should be updated with empty pile")
	}
	if postState.Players[current].CardsInHand != 4 {
		t.Errorf("hands' state should be updated (%d != 4)",
			postState.Players[current].CardsInHand)
	}

	if len(session.history) != 1 {
		t.Error("commands should be kept in history")
	}

	invalidCmd := NewThrowCommand(current+1, card)
	card, err = session.Dispatch(invalidCmd)
	if err == nil {
		t.Error("error should be propagated by Dispatch")
	}
	if len(session.history) != 1 {
		t.Error("invalid command should not be put to history")
	}
}
