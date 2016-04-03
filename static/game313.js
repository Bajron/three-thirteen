
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
            getPlayerDataFromQuery();
			setUpPlayers(d);
			if (myPlayer != -1) {
                updateMyHand().done(function() { continueAfterSetUp(d); });
			} else {
                continueAfterSetUp(d);
            }
        });
    });
}

function continueAfterSetUp(d) {
    updatePlayers(d.Game);
    if (myPlayer === d.Game.CurrentPlayer) {
        setUpMyMoves(d.Game);
    }
}

function setTranslations(data) {
    var assimilate;
    console.log('setting translations');
    ct = data.Data;

    ct.Consts = {};
    assimilate = function(a) {
        var k;
        for (k in a) {
            ct.Consts[a[k]] = parseInt(k);
        }
    };
    (function() {
        var toAssimilate = ['GameStates', 'PlayerStates', 'GameCommands', 'PlayerCommands'];
        var k;
        for (k in toAssimilate) {
            assimilate(ct[toAssimilate[k]]);
        }
    })();
    console.log(ct.Consts);
}

function CV(k) { return ct.Consts[k]; }

function getPlayerDataFromQuery() {
    var qm,u,v,e,p;
    u = window.location.href;
    qm = u.indexOf('?');
    v = u.slice(qm+1).split('&');
    for (e in v) {
        p = v[e].split('=');
        if (p[0] == 'name') {
            console.log('setting name: ' + p[1]);
            myName = p[1]
        }
    }
}

function setPrompt(txt) {
    console.log(txt);
    $('#my-player .prompt').text(txt);
}

function setUpMyMoves(game) {
    var p;
    p = game.Players[myPlayer];
    console.log('setting up moves ' + p.State);
    if (p.State === CV('TAKE')) {
        setPrompt('Take a card');
        $('#tt-pile .card,#tt-deck .card').draggable({
            revert: true,
            stack: '.card',
            connectoToSortable: '#my-player .hand',
            start: function() {
                $('#my-player').addClass('drop-possible');
            },
            stop: function() {
                $('#my-player').removeClass('drop-possible');
            },
        });
        $('#my-player').droppable({
            accept: '.card',
            drop: function (event, ui) {
                var d = ui.draggable,w,m,confirm;
                if (d.parent().attr('id') === 'tt-pile') {
                    m = 'TAKE_FROM_PILE';
                } else {
                    m = 'TAKE_FROM_DECK';
                }
                d.draggable('destroy');
                d.removeAttr('style');
                d.detach();
                w = $('<li class="dense"/>');
                w.append(d);
                w.appendTo('#my-player .hand');
                console.log('about to notify server about the move');
                confirm = $('#my-player .hand li').last();
                $.ajax({
                    url: '/3-13/test/' + myName + '/?move=' + CV(m),
                    dataType: 'json',
                }).done(function(data, status) {
                    if (data.Status !== 0) {
                        alert(data.Info);
                        return; // TODO reload?
                    }
                    confirm.html(cardToHtml(data.Data));
                    $('#tt-pile .card,#tt-deck .card').draggable('destroy');
                });
            }
        });
    } else if (p.State === CV('THROW')) {
        setPrompt('Throw a card back on pile');
        // allow pile to accept draggable
        $('.hand').sortable('option', {
            'start': function() {
                $('#tt-pile').addClass('drop-possible');
            },
            'stop': function() {
                $('#tt-pile').removeClass('drop-possible');
            },
        });

        $('#tt-pile').droppable({
            drop: function (event, ui) {
                var d = ui.draggable;
                var c = d.data('card');
                $('#tt-pile .card').html(d.find('.card'));
                d.remove();
                $.ajax({
                    'url': '/3-13/test/'+ myName+'/?move=' + CV('THROW') + '&card=' + cardToAscii(c),
                    'dataType': 'json',
                }).done(function(data, status) {
                    if (data.Status !== 0) {
                        alert(data.Info);
                        return; // TODO reload?
                    }
                    // TODO some event queue to reload and continue?
                });
            },
        });
    } else if (p.State == CV('DONE')) {
        setPrompt('Set up groups or pass the turn');
        // enable done button
        // enable setup groups button
    } else {
        setPrompt(p.State);
    }
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
	t = makeMeLast(myName, d);
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
        if (myPlayer == i) continue;
		updateCardsInHand(game.Players[i]);
	}
}

function updateCardsInHand(player) {
	var c,p,l;
	c = player.CardsInHand;
	p = $(pId(player.Index));
	l = p.find('ul.fan li').length;
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
	var h, f, i;

    $('#my-player').append(
            playerHtml(myName),
            '<div class="prompt"></div>');
	h = $('#my-player .player-box');
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
        f.find('li').last().data('card', hand[i]);
	}
}

function updateMyHand() {
	return $.ajax({
		url: '/3-13/test/'+myName+'/?hand',
		dataType: 'json',
	}).done(function(data, status) {
		addMyPlayer(data.Data);
	});
}


