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

func TestCreate104Deck(t *testing.T) {
	deck := Create104Deck()
	if len(deck) != 104 {
		t.Errorf("deck has %d cards instead of 104", len(deck));
	}

	var idx int
	if idx = deck.Find(Card{ACE, SPADES}); -1 == idx {
		t.Error("Ace of spades not found")
	}

	cutOnAce := deck[idx+1:]
	if idx = cutOnAce.Find(Card{ACE, SPADES}); -1 == idx {
		t.Error("second Ace of spades not found")
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

func TestDraw(t *testing.T) {
	deck := Create52Deck()
	_ = deck.Draw()
	if 51 != len(deck) {
		t.Error("draw did not reduce length")
	}
	
	deck = Deck(make([]Card, 0))
	card := deck.Draw()
	if 0 != len(deck) {
		t.Error("draw from empty changed length")
	}
	if card != NIL_CARD {
		t.Error("got a card from empty deck")
	}
}

func TestPop(t *testing.T) {
	pile := Pile(make([]Card, 0))
	if NIL_CARD != pile.Pop() {
		t.Error("you should get NIL_CARD from empty pile")
	}

	kh := Card{KING, HEARTS}
	as := Card{ACE, SPADES}
	pile = append(pile, as, kh)
	if kh != pile.Pop() {
		t.Error("you should get last appended with Pop")
	}
	if len(pile) != 1 {
		t.Error("Pop should remove element from the pile")
	}
	if as != pile.Top() {
		t.Error("you should see top element from a pile")
	}
}

