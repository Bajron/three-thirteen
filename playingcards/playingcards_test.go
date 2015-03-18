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
	verifySuitString(t, "S", Suit(SPADES));
	verifySuitString(t, "D", Suit(DIAMONDS));
	verifySuitString(t, "C", Suit(CLUBS));
	verifySuitString(t, "H", Suit(HEARTS));
	verifySuitString(t, "?", Suit(-1));
}

