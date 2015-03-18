package playingcards

import "testing"

func TestCreate52Deck(t *testing.T) {
	deck := Create52Deck()
	if len(deck) != 52 {
		t.Errorf("deck has %d cards instead of 52", len(deck))
	}

	if -1 == deck.Find(Card{ACE, SPADES}) {
		t.Error("Ace of spades not found")
	}
}

func TestFind(t *testing.T) {
	deck := Deck(make([]Card, 0))
	if -1 != deck.Find(Card{ACE, SPADES}) {
		t.Error("found in empty deck")
	}

	deck = Deck(make([]Card, 1))
	deck[0] = Card{ACE, SPADES}
	
	if 0 != deck.Find(Card{ACE, SPADES}) {
		t.Error("not found inserted card")
	}
}
