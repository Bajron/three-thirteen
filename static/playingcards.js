
function cardHtml(c) {
    var str;
    str =  '<span class="card"><span class="cwrap">';
    str += c;
    str += '</span></span>';
    return str;
}

function cardsEqual(a,b) {
	if (typeof a === 'undefined') {
		console.log('*!* lhs in cardsEqual is undefined');
		return false;
	}
	if (typeof b === 'undefined') {
		console.log('*!* rhs in cardsEqual is undefined');
		return false;
	}
	if (a == b) return true;
    return a.Rank === b.Rank && a.Suit === b.Suit;
}

function cardBackHtml() {
	return cardHtml('🂠');
}
function cardAsciiToUc(a) {
	if (!a) return '';
    return ''+ ct.Cards[a];
}

function cardToAscii(obj) {
	if (typeof obj === 'undefined') return 'undefined';
	if (obj.Rank < 0 || obj.Suit < 0) return '';
    return ''+ ct.Ranks[obj.Rank] + ct.Suits[obj.Suit];
}

function cardToHtml(obj) {
	return cardHtml(cardAsciiToUc(cardToAscii(obj)));
}


