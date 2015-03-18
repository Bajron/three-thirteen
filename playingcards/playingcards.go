// Copyright 2015 Blazej Kramarski. All rights reserved.

// Package playingcards provides types and functions for card games.

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
	JOCKER = 0
	ACE = 1 
	JACK = 11 + iota
	QUEEN
	KING
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
		default: return fmt.Sprintf("?%d?", int8(r))
	}
}

func (c Card) String() string {
	return fmt.Sprintf("%s %d", c.Rank, c.Suit)
}

