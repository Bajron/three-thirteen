
/* global myPlayer */
/* global myName */
/* global ct */

function initialSequence() {
    $.ajax({
        url: "/3-13/?translations",
        dataType: 'json',
        success: setTranslations
    }).done(function() {
        $.ajax({
            url: "/3-13/test/",
            dataType: 'json',
        }).done(function(data, status) {
            var d;
			d = data.Data;
            $('#tt-pile').html(cardToHtml(d.Game.PileTop));
            $('#tt-deck').html(cardBackHtml());
            console.log('got game status');
			extendSessionData(d);
			setUpPlayers(d);
			updatePlayers(d.Game);

			if (myPlayer == -1) {
				return;
			}
			updateMyHand();
        });
    });
}



function setTranslations(data) {
    console.log('setting translations');
    ct = data.Data;
}

function playerHtml(name) {
    return [
        '<div class="player-box">',
        '<span class="player-name">'+ name +'</span>',
        '<ul class="fan"></ul>',
        '</div>'
    ].join('');
}

/* WARNING: inefficient IMHO */
function rotateArr(a, i) {
	return a.concat(a).slice(i, i + a.length);
}

function extendSessionData(d) {
	var i;
	for (i=0; i<d.Players.length;++i) {
		d.Game.Players[i].Index = i;
		d.Game.Players[i].Name = d.Players[i];
	}
}

function setUpPlayers(d) {
	var i,t;
	t = makeMeLast(d.Players[0], d);
	if (t.Me != -1) {
		t.Names.pop();
		t.Players.pop();
		myPlayer = t.Me;
	}
	addPlayers(t.Players);
}


function addPlayers(players) {
	var i,j,h,f,p,target;
	target = $('#other-players');
	for (i=0;i<players.length;++i) {
		p = players[i];
		h = $(playerHtml(p.Name));
		h.attr('id', 'player-' + p.Index);
		f = h.find('ul.fan');
		for (j=0; j<p.CardsInHand;++j) {
			f.append('<li>'+cardBackHtml()+'</li>');
		}
		target.append(h);
	}
}

function makeMeLast(name, d) {
	var me, rot;
	me = d.Players.indexOf(name);
	if (me != -1) {
		rot = (me + 1) % d.Players.length;
		return {
			'Me': me,
			'Names': rotateArr(d.Players, rot),
			'Players':  rotateArr(d.Game.Players, rot),
		};
	} else {
		return {
			'Me': me,
			'Names': d.Players,
			'Players':  d.Game.Players,
		};
	}
}

function pId(i) {
	return '#player-' + i;
}

function updatePlayers(game) {
	var i;
	$('.active-player').removeClass('active-player');
	$(pId(game.CurrentPlayer)).addClass('active-player');

	for (i=0; i<game.Players.length; ++i) {
		updateCardsInHand(game.Players[i]);
	}
}

function updateCardsInHand(player) {
	var c,p,l;
	c = player.CardsInHand;
	p = $(pId(player.Index));
	l = p.find('ul li').length;
	while (l > c) {
		p.find('ul li').last().remove();
		--l;
	}
	while (l < c) {
		p.find('ul').append('<li>'+cardBackHtml()+'</li>');
		++l;
	}
}

function addMyPlayer(hand) {
	var target, h, f, i;

	target = $('#my-player');

	h = $(playerHtml(myName));
	h.attr('id', 'player-' + myPlayer);

	f = h.find('ul.fan');
	f.removeClass('fan');
	f.addClass('hand');
	f.wrap('<div class="hwrap"/>');
	f.sortable({
    	placeholder: 'hand-placeholder',
        revert: 250,
    });
    f.disableSelection();
    
	for (i=0; i<hand.length;++i) {
		f.append('<li class="dense">'+ cardToHtml(hand[i]) +'</li>');
	}
	target.append(h);
}

function updateMyHand() {
	$.ajax({
		url: '/3-13/test/'+myName+'/?hand',
		dataType: 'json',
	}).done(function(data, status) {
		addMyPlayer(data.Data);
	});
}


