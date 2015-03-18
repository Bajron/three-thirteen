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
	verifyRankString(t, "I", Rank(JOCKER))
	verifyRankString(t, "J", Rank(JACK))
	verifyRankString(t, "Q", Rank(QUEEN))
	verifyRankString(t, "K", Rank(KING))
	verifyRankString(t, "10", Rank(10))
	verifyRankString(t, "2", Rank(2))
	verifyRankString(t, "?-1?", Rank(-1))
}
