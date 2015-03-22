package playingcards

import (
	"sort"
)

type Hand []Card
type Group []Card
type rankMatch func (Rank) bool

func (g Group) Len() int {
	return len(g)
}

func (g Group) Less(i, j int) bool {
	return Compare(g[i], g[j]) < 0;
}

func (g Group) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func IsSet(g Group, isJocker rankMatch) bool {
	if len(g) < 3 {
		return false
	}

	rank, i := firstNotJockerRank(g, isJocker)

	for _, c := range g[i:] {
		if !isJocker(c.Rank) && rank != c.Rank {
			return false
		}
	}
	return true
}

func firstNotJockerRank(g Group, isJocker rankMatch) (Rank, int) {
	var rank Rank
	i := 0
	for ; i<len(g) && isJocker(g[i].Rank); i++ {}
	if i < len(g) {
		rank = g[i].Rank
	}
	return rank, i
}

func NewGroup(c ...Card) (g Group) {
	g = append(g, c...)
	return
}

func IsSequence(g Group, isJocker rankMatch) bool {
	if len(g) < 3 {
		return false
	}

	firstNotJocker := sortPutJockersFirst(g, isJocker)

	if firstNotJocker >= len(g) {
		return true
	}

	i := firstNotJocker
	s := g[i].Suit
	r := g[i].Rank
	jockersToUse := i

	for _, c := range g[i:] {
		if c.Suit != s {
			return false
		}
		for ;jockersToUse > 0 && c.Rank != r ; {
			jockersToUse--
			r++
		}
		if c.Rank != r {
			return false
		}
		r++
	}

	return true
}

func sortPutJockersFirst(g Group, isJocker rankMatch) (firstNotJocker int) {
	sort.Sort(g)

	var i int = 0
	for ; i<len(g) && isJocker(g[i].Rank); i++ {}

	firstNotJocker = i
	for i = firstNotJocker + 1; i<len(g); i++ {
		if isJocker(g[i].Rank) {
			g.Swap(firstNotJocker, i)
			firstNotJocker++
		}
	}

	return firstNotJocker
}

