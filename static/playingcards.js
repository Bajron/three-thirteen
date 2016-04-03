
function cardHtml(c) {
    var str;
    str =  '<span class="card"><span class="cwrap">';
    str += c;
    str += '</span></span>';
    return str;
}

function cardBackHtml() {
	return cardHtml('🂠');
}
function cardAsciiToUc(a) {
	if (!a) return '';
    return ''+ ct.Cards[a];
}

function cardToAscii(obj) {
	if (obj.Rank < 0 || obj.Suit < 0) return '';
    return ''+ ct.Ranks[obj.Rank] + ct.Suits[obj.Suit];
}

function cardToHtml(obj) {
	return cardHtml(cardAsciiToUc(cardToAscii(obj)));
}


