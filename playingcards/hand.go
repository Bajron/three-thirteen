package playingcards

type Hand []Card
type Group []Card
type rankMatch func (Rank) bool

func IsSet(g Group, isJocker rankMatch) bool {
	if len(g) < 3 {
		return false
	}

	var rank Rank
	i := 0
	for ; i<len(g) && isJocker(g[i].Rank); i++ {}
	if i < len(g) {
		rank = g[i].Rank
	}

	for _, c := range g[i:] {
		if !isJocker(c.Rank) && rank != c.Rank {
			return false
		}
	}
	return true
}

func NewGroup(c ...Card) (g Group) {
	g = append(g, c...)
	return
}

