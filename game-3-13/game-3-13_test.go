package game313

import (
	"testing"
	"github.com/Bajron/three-thirteen/playingcards"
)

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

func TestMove(t *testing.T) {
	state := New(3)
	state.Deal()

	card, err := state.TakeMove(2, TAKE_FROM_PILE)
	if card != playingcards.NIL_CARD || err == nil {
		t.Error("move should fail if not current player")
	}
	card, err = state.TakeMove(0, -1)
	if card != playingcards.NIL_CARD || err == nil {
		t.Error("move should fail if invalid move")
	}

	onPile := state.Pile.Top()
	card, err = state.TakeMove(0, TAKE_FROM_PILE)
	if card != onPile || err != nil {
		t.Error("move should return a card from the pile")
	}

	if state.Players[0].Hand[len(state.Players[0].Hand) - 1] != onPile {
		t.Error("player's hand should be updated")
	}

	deckLen := len(state.Deck)
	card, err = state.TakeMove(0, TAKE_FROM_DECK)
	if card != playingcards.NIL_CARD ||
			err == nil || deckLen != len(state.Deck) {
		t.Error("second take move should fail")
	}

	// reinitialized
	state = New(3)
	state.Deal()

	deckLen = len(state.Deck)
	card, err = state.TakeMove(0, TAKE_FROM_DECK)
	if card == playingcards.NIL_CARD ||
			err != nil || (deckLen - 1) != len(state.Deck) {
		t.Error("move should take a card from deck")
	}

	// TODO: throw move
	// TODO: done move
}

