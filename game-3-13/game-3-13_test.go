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

	gotCard := card
	deckLen := len(state.Deck)
	card, err = state.TakeMove(0, TAKE_FROM_DECK)
	if card != playingcards.NIL_CARD ||
		err == nil || deckLen != len(state.Deck) {
		t.Error("second take move should fail")
	}

	err = state.ThrowMove(0, gotCard)
	if err != nil {
		t.Errorf("should be able to throw taken card (%v)", err)
	}
	card, err = state.TakeMove(0, TAKE_FROM_DECK)
	if card != playingcards.NIL_CARD ||
		err == nil || deckLen != len(state.Deck) {
		t.Error("take move after throw should fail")
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
	state.PassMove(0)
	if state.CurrentPlayer != 1 {
		t.Error("current player should be 1")
	}
	card, _ = state.TakeMove(1, TAKE_FROM_PILE)
	state.ThrowMove(1, card)
	state.PassMove(1)
	if state.CurrentState != STARTING {
		t.Error("state should be still STARTING")
	}
	card, _ = state.TakeMove(2, TAKE_FROM_PILE)
	state.ThrowMove(2, card)
	state.PassMove(2)
	if state.CurrentState != PLAYING {
		t.Error("state should be PLAYING")
	}
	if state.CurrentPlayer != 0 {
		t.Error("current player should be 0")
	}
}

func TestDoneMove(t *testing.T) {
	validDoneHand := playingcards.Group{
		{5, playingcards.HEARTS},
		{6, playingcards.HEARTS},
		{7, playingcards.HEARTS},
	}
	// remember 3 is the jocker
	invalidDoneHand := playingcards.Group{
		{4, playingcards.HEARTS},
		{6, playingcards.HEARTS},
		{7, playingcards.HEARTS},
	}
	state := prepareStateWithHand(0, 3, validDoneHand)
	takeAndThrow(state, 0)

	err := state.DoneMove(0, makeSingleGroupSetFromHand(validDoneHand))
	if err == nil {
		t.Error("you cannot finish in first round")
	}

	state.CurrentState = PLAYING
	err = state.DoneMove(1, makeSingleGroupSetFromHand(validDoneHand))
	if err == nil {
		t.Error("you cannot finish not on your turn")
	}

	finalGroups := makeSingleGroupSetFromHand(validDoneHand)
	finalGroups.Unassigned = playingcards.Group{{6, playingcards.CLUBS}}
	err = state.DoneMove(0, finalGroups)
	if err == nil {
		t.Error("you should not finish with unassigned cards")
	}

	// REINITIALIZATION
	state = prepareStateWithHand(0, 3, invalidDoneHand)
	state.CurrentState = PLAYING

	err = state.DoneMove(0, makeSingleGroupSetFromHand(validDoneHand))
	if err == nil {
		t.Errorf("you should finish only after take and throw")
	}

	takeAndThrow(state, 0)

	err = state.DoneMove(0, makeSingleGroupSetFromHand(invalidDoneHand))
	if err == nil {
		t.Error("you should not finish with invalid set")
	}

	// REINITIALIZATION
	state = prepareStateWithHand(0, 3, validDoneHand)
	state.CurrentState = PLAYING
	takeAndThrow(state, 0)

	validDoneHandCheat := playingcards.Group{
		{5, playingcards.CLUBS},
		{6, playingcards.CLUBS},
		{7, playingcards.CLUBS},
	}

	err = state.DoneMove(0, makeSingleGroupSetFromHand(validDoneHandCheat))
	if err == nil {
		t.Error("cheat moves should be detected and rejected")
	}

	err = state.DoneMove(0, makeSingleGroupSetFromHand(validDoneHand))
	if err != nil {
		t.Errorf("valid move should be possible '%v'", err)
	}
	if state.CurrentState != FINISHING {
		t.Error("after first done, we should go to final round")
	}
	if state.CurrentPlayer != 1 {
		t.Error("current player should be 1")
	}

	// completion
	takeAndThrow(state, 1)
	err = state.DoneMove(1, makeAllUnassigned(state, 1))
	assertNoError(t, err)
	takeAndThrow(state, 2)
	err = state.DoneMove(2, makeAllUnassigned(state, 2))
	assertNoError(t, err)

	if state.CurrentState != FINISHED {
		t.Error("game should be finished after final round")
	}

	_, err = state.TakeMove(0, TAKE_FROM_PILE)
	if err == nil {
		t.Error("no moves should be possible after finishing round")
	}

	state.AdvanceRound()
	if state.Players[0].Points != 0 {
		t.Error("player with no unassigned should have 0 points")
	}
	if state.Players[1].Points == 0 {
		t.Error("player with unassigned should not have 0 points")
	}
	if state.Trumph != 4 {
		t.Error("trumph should be advanced")
	}
	for i := 0; i < 3; i++ {
		if len(state.Players[i].Hand) > 0 {
			t.Error("hands should be cleared")
		}
		if len(state.FinalGroups[i].Set) > 0 || len(state.FinalGroups[i].Unassigned) > 0 {
			t.Error("final groups should be cleared")
		}
	}
	if len(state.Pile) != 0 {
		t.Error("pile should be empty")
	}
	if len(state.Deck) != 104 {
		t.Error("deck should be reset")
	}
}

func TestAdvanceRound(t *testing.T) {
	state := newWithStaringPlayer(3, 0)
	state.Deal()

	err := state.AdvanceRound()
	if err == nil {
		t.Error("you should not advance round when it is not finished")
	}

	state.CurrentState = FINISHED
	state.Trumph = 14
	err = state.AdvanceRound()
	if state.Trumph != 14 {
		t.Error("you should not advance after trumph 14")
	}
	if len(state.Deck) != 0 || len(state.Pile) != 0 {
		t.Error("deck and pile should be removed after t14")
	}
}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Errorf("unexpected error: %v", e)
	}
}

func TestPass(t *testing.T) {
	state := newWithStaringPlayer(3, 0)
	state.Deal()

	err := state.PassMove(0)
	if err == nil {
		// TODO still can be simulated by take and throw from pile
		t.Error("you should not end turn without any moves")
	}

	card, _ := state.TakeMove(0, TAKE_FROM_PILE)

	err = state.PassMove(0)
	if err == nil {
		t.Error("you should not end turn without any moves")
	}

	state.ThrowMove(0, card)

	err = state.PassMove(1)
	if err == nil {
		t.Error("you should not act as different player")
	}

	state.CurrentState = FINISHING
	err = state.PassMove(0)
	if err == nil {
		t.Error("you should only do Done move on final round")
	}

	state.CurrentState = PLAYING

	err = state.PassMove(0)
	if err != nil {
		t.Error("you should be able to succeed with PassMove")
	}
	if state.Players[0].State != WAIT {
		t.Error("player state should be default")
	}
	if state.CurrentPlayer != 1 {
		t.Error("current player should be advanced")
	}
}

func prepareStateWithHand(player int, players int, hand playingcards.Group) *State {
	state := newWithStaringPlayer(players, player)
	prepareDeck(state.Deck, player, players, hand)
	state.Deal()
	return state
}

func takeAndThrow(state *State, player int) {
	card, _ := state.TakeMove(player, TAKE_FROM_PILE)
	state.ThrowMove(player, card)
}

func prepareDeck(d playingcards.Deck, player int, players int, hand playingcards.Group) playingcards.Deck {
	for i, c := range hand {
		from := player + i*players
		for j := from + 1; j < len(d); j++ {
			if d[j] == c {
				tmp := d[j]
				d[j] = d[from]
				d[from] = tmp
			}
		}
	}
	return d
}

func makeSingleGroupSetFromHand(h playingcards.Group) FinalGroups {
	cards := []playingcards.Card(h)
	groups := make([]playingcards.Group, 1)
	groups[0] = cards
	return FinalGroups{groups, playingcards.Group{}}
}

func makeAllUnassigned(state *State, player int) FinalGroups {
	return FinalGroups{nil, state.Players[player].Hand}
}

func (s State) curPlayerHasValidHand() bool {
	return len(*s.currentPlayerHand()) == int(s.Trumph)
}
