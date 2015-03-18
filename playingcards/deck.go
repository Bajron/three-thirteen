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

func (deck *Deck) Find(c Card) (int) {
	for i, v := range *deck {
		if v == c {
			return i
		}
	}
	return -1
}

