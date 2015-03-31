package game313

import "testing"

func TestNew(t *testing.T) {
	var value interface{}
	value = New(3)
	state, ok := value.(*State)
	if ! ok {
		t.Error("returned value is not of type State*")
	}

	if len(state.Players) != 3 {
		t.Error("new game does not have given number of players")
	}

	if len(state.Deck) != 104 {
		t.Error("new game should be played with 104 cards")
	}

	if len(state.Pile) != 0 {
		t.Error("new game should have no cards on the pile")
	}

	if state.Trumph != 3 {
		t.Error("new game should start with 3 as trumph")
	}

	if len(state.Players[0].Hand) != 0 {
		t.Error("initial hand should be empty")
	}

	if state.Players[0].Points != 0 {
		t.Error("initial points should be 0")
	}
}

func TestDeal(t *testing.T) {
	state := New(3)
	state.Deal()

	if l := len(state.Players[0].Hand); l != 3 {
		t.Errorf("initial hand should have 3 cards (has %d)", l)
	}

	if len(state.Pile) != 1 {
		t.Error("there should be one card on the pile initially")
	}

	if len(state.Deck) != (104 - 3*3 - 1) {
		t.Error("deck should not have already dealed cards")
	}
}

