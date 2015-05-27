package game313

import (
	"github.com/Bajron/three-thirteen/playingcards"
	"testing"
)

func TestNew(t *testing.T) {
	var value interface{}
	value = New(3)
	state, ok := value.(*State)
	if !ok {
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
	state := newWithStaringPlayer(3, 0)
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

	if state.Players[0].Hand[len(state.Players[0].Hand)-1] != onPile {
		t.Error("player's hand should be updated")
	}

	deckLen := len(state.Deck)
	card, err = state.TakeMove(0, TAKE_FROM_DECK)
	if card != playingcards.NIL_CARD ||
		err == nil || deckLen != len(state.Deck) {
		t.Error("second take move should fail")
	}

	// reinitialized
	state = newWithStaringPlayer(3, 0)
	state.Deal()

	deckLen = len(state.Deck)
	card, err = state.TakeMove(0, TAKE_FROM_DECK)
	if card == playingcards.NIL_CARD ||
		err != nil || (deckLen-1) != len(state.Deck) {
		t.Error("move should take a card from deck")
	}

	// TODO: done move
}

func TestThrowMove(t *testing.T) {
	state := newWithStaringPlayer(3, 0)
	state.Deal()

	err := state.ThrowMove(0, (*state.currentPlayerHand())[0])
	if err == nil {
		t.Error("throw without earlier take is error")
	}

	err = state.ThrowMove(1, (*state.currentPlayerHand())[0])
	if err == nil {
		t.Error("throw of different player is error")
	}

	// ready to throw after take
	state.TakeMove(0, TAKE_FROM_DECK)

	err = state.ThrowMove(1, (*state.currentPlayerHand())[0])
	if err == nil {
		t.Error("throw of different player is error")
	}

	err = state.ThrowMove(0, playingcards.NIL_CARD)
	if err == nil {
		t.Error("throwing not owned card should fail")
	}

	card := (*state.currentPlayerHand())[0]
	err = state.ThrowMove(0, card)
	if err != nil {
		t.Error("throw should succeed")
	}
	if !state.curPlayerHasValidHand() {
		t.Error("hand should have correct number of cards")
	}
	if state.Pile.Top() != card {
		t.Error("thrown card should be on top of pile")
	}
}

func TestStartingToPlayingTransition(t *testing.T) {
	state := newWithStaringPlayer(3, 0)
	state.Deal()

	card, _ := state.TakeMove(0, TAKE_FROM_PILE)
	state.ThrowMove(0, card)
	if state.CurrentPlayer != 1 {
		t.Error("current player should be 1")
	}
	card, _ = state.TakeMove(1, TAKE_FROM_PILE)
	state.ThrowMove(1, card)
	if state.CurrentState != STARTING {
		t.Error("state should be still STARTING")
	}
	card, _ = state.TakeMove(2, TAKE_FROM_PILE)
	state.ThrowMove(2, card)
	if state.CurrentState != PLAYING {
		t.Error("state should be PLAYING")
	}
	if state.CurrentPlayer != 0 {
		t.Error("current player should be 0")
	}
}

func TestDoneMove(t *testing.T) {
	validDoneHand := playingcards.Hand{
		{5, playingcards.HEARTS},
		{6, playingcards.HEARTS},
		{7, playingcards.HEARTS},
	}
	state := newWithStaringPlayer(3, 0)
	state.Deal()

	state.Players[0].Hand = validDoneHand

	err := state.DoneMove(0, makeSingleGroupSetFromHand(validDoneHand))
	if err == nil {
		t.Error("you cannot finish in first round")
	}

	state.CurrentState = PLAYING

	err = state.DoneMove(1, makeSingleGroupSetFromHand(validDoneHand))
	if err == nil {
		t.Error("you cannot finish not on your turn")
	}

	// TODO first done must have not unassigned cards
	// TODO completion round
	// TODO points
	// TODO advance to the next trumph game
}

func makeSingleGroupSetFromHand(h playingcards.Hand) FinalGroups {
	cards := []playingcards.Card(h)
	groups := make([]playingcards.Group, 1)
	groups[0] = cards
	return FinalGroups{groups, playingcards.Group{}}
}

func (s State) curPlayerHasValidHand() bool {
	return len(*s.currentPlayerHand()) == int(s.Trumph)
}
