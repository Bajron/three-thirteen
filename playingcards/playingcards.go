// Copyright 2015 Blazej Kramarski. All rights reserved.

// Package playingcards provides types and functions for card games.
package playingcards

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

var SUITS []Suit = nil
var RANKS []Rank = nil
var NIL_CARD = Card{-1, -1}

func init() {
	SUITS = make([]Suit, SUIT_COUNT)
	for i := uint(1); (1 << i) != SUIT_GUARD; i++ {
		SUITS[i] = Suit(1 << i)
	}
	RANKS = make([]Rank, RANK_COUNT)
	for i := 0; i < int(RANK_COUNT); i++ {
		RANKS[i] = Rank(i)
	}
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// UString returns unicode rune for a card
// http://en.wikipedia.org/wiki/Playing_cards_in_Unicode
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

func Compare(lhs Card, rhs Card) int {
	if lhs.Rank < rhs.Rank ||
		lhs.Rank == rhs.Rank && lhs.Suit < rhs.Suit {
		return -1
	}
	if lhs.Rank > rhs.Rank ||
		lhs.Rank == rhs.Rank && lhs.Suit > rhs.Suit {
		return 1
	}
	return 0
}

type Rank int8

const (
	JACK = 11 + iota
	QUEEN
	KING
	RANK_COUNT
	JOCKER = 0
	ACE    = 1
)

func (r Rank) String() string {
	switch {
	case 2 <= r && r <= 10:
		return strconv.Itoa(int(r))
	case r == ACE:
		return "A"
	case r == JACK:
		return "J"
	case r == QUEEN:
		return "Q"
	case r == KING:
		return "K"
	case r == JOCKER:
		return "*"
	default:
		return fmt.Sprintf("?%d?", int8(r))
	}
}

type Suit int8

const (
	CLUBS = 1 << iota
	DIAMONDS
	SPADES
	HEARTS
	SUIT_GUARD
	SUIT_COUNT = 4
)

func getBaseUnicodeCard(s Suit) rune {
	switch s {
	case CLUBS:
		return 0x1F0D0
	case DIAMONDS:
		return 0x1F0C0
	case SPADES:
		return 0x1F0A0
	case HEARTS:
		return 0x1F0B0
	default:
		return utf8.RuneError
	}
}

func (s Suit) String() string {
	switch s {
	case CLUBS:
		return "C"
	case DIAMONDS:
		return "D"
	case SPADES:
		return "S"
	case HEARTS:
		return "H"
	default:
		return "?"
	}
}

func (s Suit) UString() string {
	switch s {
	case CLUBS:
		return "\u2664"
	case DIAMONDS:
		return "\u2666"
	case SPADES:
		return "\u2660"
	case HEARTS:
		return "\u2665"
	default:
		return "?"
	}
}
