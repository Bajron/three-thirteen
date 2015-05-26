package playingcards

import "testing"

func TestIsSet(t *testing.T) {
	var h Group
	jocker := func(r Rank) bool { return r == JOCKER }
	h = NewGroup()
	if IsSet(h, jocker) {
		t.Error("empty group should not be set")
	}

	h = NewGroup(Card{}, Card{})
	if IsSet(h, jocker) {
		t.Error("group smaller than 3 cannot be set")
	}

	h = NewGroup(Card{3, 2}, Card{3, 2}, Card{2, 2})
	if IsSet(h, jocker) {
		t.Error("different ranks are not set")
	}

	h = NewGroup(Card{JOCKER, 2}, Card{2, 2}, Card{2, 2})
	if !IsSet(h, jocker) {
		t.Error("2 same and jocker should be set")
	}

	h = NewGroup(Card{JOCKER, 2}, Card{JOCKER, 2}, Card{2, 2})
	if !IsSet(h, jocker) {
		t.Error("1 and 2*jocker should be set")
	}

	h = NewGroup(Card{JOCKER, 2}, Card{2, 2}, Card{JOCKER, 2})
	if !IsSet(h, jocker) {
		t.Error("2 jockers and other should be set")
	}
}

func TestIsSequence(t *testing.T) {
	var h Group
	jocker := func(r Rank) bool { return r == JOCKER }

	h = NewGroup()
	if IsSequence(h, jocker) {
		t.Error("empty set is not a sequence")
	}

	h = NewGroup(Card{}, Card{})
	if IsSequence(h, jocker) {
		t.Error("group smaller than 3 cannot be sequence")
	}

	h = NewGroup(Card{3, HEARTS}, Card{4, HEARTS}, Card{5, HEARTS})
	if !IsSequence(h, jocker) {
		t.Error("3 cards in sequence should be a sequence")
	}

	h = NewGroup(Card{3, HEARTS}, Card{4, HEARTS}, Card{5, CLUBS})
	if IsSequence(h, jocker) {
		t.Error("different suits are not set")
	}

	h = NewGroup(Card{3, HEARTS}, Card{4, HEARTS}, Card{JOCKER, CLUBS})
	if !IsSequence(h, jocker) {
		t.Error("2 in sequence and a jocker should be a sequence")
	}

	j := Card{JOCKER, HEARTS}
	h = NewGroup(j, j, j)
	if !IsSequence(h, jocker) {
		t.Error("3 jockers should be a sequence")
	}

	h = NewGroup(Card{3, HEARTS}, Card{5, HEARTS}, Card{4, HEARTS})
	if !IsSequence(h, jocker) {
		t.Error("3 cards in sequence shuffled should be a sequence")
	}

	h = NewGroup(Card{3, HEARTS}, Card{4, HEARTS}, Card{6, HEARTS})
	if IsSequence(h, jocker) {
		t.Error("3 cards in sequence with a gap " +
			"should not be a sequence")
	}

	h = NewGroup(Card{JOCKER, CLUBS}, Card{3, HEARTS}, Card{4, HEARTS})
	if !IsSequence(h, jocker) {
		t.Error("jocker and 2 in sequence should be a sequence")
	}

	h = NewGroup(Card{JOCKER, CLUBS}, Card{2, HEARTS}, Card{4, HEARTS})
	if !IsSequence(h, jocker) {
		t.Error("jocker should fill gaps in sequence")
	}

	h = NewGroup(Card{JOCKER, CLUBS}, Card{2, HEARTS}, Card{4, HEARTS}, Card{5, HEARTS})
	if !IsSequence(h, jocker) {
		t.Error("jocker should fill gaps in sequence")
	}

	fourIsJocker := func(r Rank) bool { return r == 4 }
	h = NewGroup(Card{3, HEARTS}, Card{4, HEARTS}, Card{4, HEARTS}, Card{6, HEARTS})
	if !IsSequence(h, fourIsJocker) {
		t.Error("jocker value is more important than nominal")
	}

	nineIsJocker := func(r Rank) bool { return r == 9 }
	h = NewGroup(Card{3, HEARTS}, Card{9, HEARTS}, Card{9, HEARTS}, Card{6, HEARTS})
	if !IsSequence(h, nineIsJocker) {
		t.Error("jocker value is more important than nominal in the end")
	}
}
