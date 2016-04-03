package playingcards

import "testing"

func verifyRankString(t *testing.T, s string, r Rank) {
	if s != r.String() {
		t.Errorf("%s != %s", s, r)
	}
}

func TestRankString(t *testing.T) {
	verifyRankString(t, "A", Rank(1))
	verifyRankString(t, "A", Rank(ACE))
	verifyRankString(t, "*", Rank(JOCKER))
	verifyRankString(t, "J", Rank(JACK))
	verifyRankString(t, "Q", Rank(QUEEN))
	verifyRankString(t, "K", Rank(KING))
	verifyRankString(t, "10", Rank(10))
	verifyRankString(t, "2", Rank(2))
	verifyRankString(t, "?-1?", Rank(-1))
}

func verifySuitString(t *testing.T, s string, su Suit) {
	if s != su.String() {
		t.Errorf("%s != %s", s, su)
	}
}

func TestSuitString(t *testing.T) {
	verifySuitString(t, "S", Suit(SPADES))
	verifySuitString(t, "D", Suit(DIAMONDS))
	verifySuitString(t, "C", Suit(CLUBS))
	verifySuitString(t, "H", Suit(HEARTS))
	verifySuitString(t, "?", Suit(-1))
}

func verifyCompare(t *testing.T, lhs Card, rhs Card, expected int) {
	r := Compare(lhs, rhs)
	if r != expected {
		t.Errorf("invalid result %d (!= %d) for %s and %s",
			r, expected, lhs, rhs)
	}
}

func TestCompare(t *testing.T) {
	verifyCompare(t, Card{2, HEARTS}, Card{2, HEARTS}, 0)
	verifyCompare(t, Card{3, HEARTS}, Card{4, HEARTS}, -1)
	verifyCompare(t, Card{3, HEARTS}, Card{2, HEARTS}, 1)
}

func TestNewFromString(t *testing.T) {
	verifyFromString(t, "", NIL_CARD)
	verifyFromString(t, "2C", Card{2, CLUBS})
	verifyFromString(t, "10C", Card{10, CLUBS})
	verifyFromString(t, "12C", NIL_CARD)
	verifyFromString(t, "QC", Card{QUEEN, CLUBS})
	verifyFromString(t, "QH", Card{QUEEN, HEARTS})
}

func verifyFromString(t *testing.T, s string, expected Card) {
	got := NewFromString(s)
	if expected != got {
		t.Errorf("invalid result %s != %s (from %s)",
			got, expected, s)
	}
}
