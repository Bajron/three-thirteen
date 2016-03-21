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
	session := NewGameSession(3)

	if len(session.state.Players) != 3 {
		t.Error("session should have initialized state")
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
	/* TODO fix dealing
	if postState.Players[current].CardsInHand != 4 {
		t.Errorf("hands' state should be updated (%d != 4)",
			postState.Players[current].CardsInHand)
	}
	*/
}
