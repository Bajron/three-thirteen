package playingcards 

import (
	"fmt"
	"strconv"
)

type Card struct {
	Rank Rank
	Suit Suit
}

type Suit int8 

const (
	CLUBS = iota
	DIAMONDS
	SPADES
	HEARTS
)

type Rank int8 

const (
	ACE = 1 
	JACK = 11
	QUEEN
	KING
	JOCKER
)

func (r Rank) String() string {
	switch {
		case 2 <= r && r <= 10 :
			return strconv.Itoa(int(r))
		case r == ACE: return "A"
		case r == JACK: return "J"
		case r == QUEEN: return "Q"
		case r == KING: return "K"
		case r == JOCKER: return "I"
		default: return fmt.Sprint("?%d?", r)
	}
}

func (c Card) String() string {
	return fmt.Sprintf("%s %d", c.Rank, c.Suit)
}
