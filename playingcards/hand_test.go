package playingcards

import "testing"

func TestIsSet(t *testing.T) {
	var h Group
	jocker := func (r Rank) bool { return r == JOCKER }
	h = NewGroup()
	if IsSet(h, jocker) {
		t.Error("empty group should not be set")
	}
	
	h = NewGroup(Card{}, Card{})
	if IsSet(h, jocker) {
		t.Error("group smaller than 3 cannot be set")
	}

	h = NewGroup(Card{3,2}, Card{3,2}, Card{2,2})
	if IsSet(h, jocker) {
		t.Error("different ranks are not set")
	}

	h = NewGroup(Card{JOCKER, 2}, Card{2, 2}, Card{2,2})
	if !IsSet(h, jocker) {
		t.Error("2 same and jocker should be set")
	}

	h = NewGroup(Card{JOCKER, 2}, Card{JOCKER, 2}, Card{2,2})
	if !IsSet(h, jocker) {
		t.Error("1 and 2*jocker should be set")
	}

	h = NewGroup(Card{JOCKER, 2}, Card{2, 2}, Card{JOCKER,2})
	if !IsSet(h, jocker) {
		t.Error("2 jockers and other should be set")
	}
}

