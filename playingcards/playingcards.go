// Copyright 2015 Blazej Kramarski. All rights reserved.

// Package playingcards provides types and functions for card games.

// http://en.wikipedia.org/wiki/Playing_cards_in_Unicode

package playingcards 

import (
	"fmt"
	"strconv"
	"unicode/utf8"
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

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

func (c Card) UString() string {
	base := getBaseUnicodeCard(c.Suit)
	if base == utf8.RuneError {
		return "???"
	}
	card := rune(int(base) + int(c.Rank))
	ret := make([]byte, utf8.RuneLen(card))
	utf8.EncodeRune(ret, card)
	return fmt.Sprint(string(ret))
}

func (c Card) StringWithUnicode() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit.UString())
}

func getBaseUnicodeCard(s Suit) rune {
	switch s {
		case CLUBS: return 0x1F0D0;
		case DIAMONDS: return 0x1F0C0;
		case SPADES: return 0x1F0A0;
		case HEARTS: return 0x1F0B0;
		default: return utf8.RuneError;
	}
}

func (s Suit) String() string {
	switch s {
		case CLUBS: return "C";
		case DIAMONDS: return "D";
		case SPADES: return "S";
		case HEARTS: return "H";
		default: return "?";
	}
}

func (s Suit) UString() string {
	switch s {
		case CLUBS: return "\u2664";
		case DIAMONDS: return "\u2666";
		case SPADES: return "\u2660";
		case HEARTS: return "\u2665";
		default: return "?";
	}
}

func (r Rank) String() string {
	switch {
		case 2 <= r && r <= 10 :
			return strconv.Itoa(int(r))
		case r == ACE: return "A"
		case r == JACK: return "J"
		case r == QUEEN: return "Q"
		case r == KING: return "K"
		case r == JOCKER: return "*"
		default: return fmt.Sprintf("?%d?", int8(r))
	}
}

