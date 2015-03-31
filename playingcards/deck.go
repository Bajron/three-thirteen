package playingcards

type Deck []Card

func Create52Deck() (Deck) {
	ret := make([]Card, 52)
	i := 0
	for _, suit := range SUITS {
		for _, rank := range RANKS[1:] {
			ret[i] = Card{rank, suit}
			i++
		}
	}
	return ret
}

func Create104Deck() (Deck) {
	ret := make([]Card, 104);
	d52 := Create52Deck()
	copy(ret, d52)
	copy(ret[52:], d52)
	return ret
}

func (deck *Deck) Find(c Card) (int) {
	for i, v := range *deck {
		if v == c {
			return i
		}
	}
	return -1
}

func (deck *Deck) Draw() (c Card) {
	if len(*deck) <= 0 {
		return NIL_CARD
	}
	c = (*deck)[0]
	*deck = (*deck)[1:]
	return
}

// TODO: create a Pile == Deck with method Pop()

