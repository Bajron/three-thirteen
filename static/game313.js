
/* global myPlayer */
/* global myName */
/* global ct */
/* global cmdQ */

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
            getPlayerDataFromQuery();
            extendSessionData(d);
            setUpPlayers(d);
            setUpTestingLinks();
            if (myPlayer != -1) {
                addMyHand().done(function() { continueAfterSetUp(d); });
            } else {
                continueAfterSetUp(d);
            }
        });
    });
}

function setUpTestingLinks() {
    console.log('TODO: remove testing links');
    $('.player-name').each(function(i,el) {
        el = $(el);
        el.wrap('<a href="?name=' + el.text() + '"/>');
    })
}

function continueAfterSetUp(d) {
    var game = d.Game;
	$('#tt-pile').html(cardToHtml(game.PileTop));
	$('#tt-pile .card').data('card', game.PileTop);

	$('#tt-deck').html(cardBackHtml());

	updatePlayers(game);

    if (game.CurrentState === CV('NOT_DEALT')) {
        if (myPlayer === game.DealingPlayer) {
		    setUpDealer();
        } else {
            setPrompt('Wait for the dealer');
        }
	} else if (myPlayer === game.CurrentPlayer) {
        setUpMyMoves(game);
    } else {
        setPrompt('Wait for your turn');
    }
    consumeCommands();
}

function setUpDealer() {
    var a;

    setPrompt('You are the dealer');

    a = $('#my-player .actions');
    a.show();
    a.find('.done,.pass').attr('disabled', 'disabled');
    a.find('.deal')
    .removeAttr('disabled')
    .click(function(ev) {
        $.ajax({
            'url': '/3-13/test/' + myName +'/?marshal=' + CV('DEAL'),
            'dataType': 'json',
        }).done(function(data, status) {
            if (data.Status !== 0) {
                alert(data.Info);
                return; // TODO reload?
            }
            cmdQ.push(synchronizeTableStatus);
        });
        ev.preventDefault();
    });
}

function consumeCommands() {
    var cmd;
    while (cmdQ.length > 0) {
        cmd = cmdQ.shift();
        cmd();
    }
    setTimeout(consumeCommands, 1000);
}

function synchronizeTableStatus() {
    $.ajax({
        'url': '/3-13/test/',
        'dataType': 'json',
    }).done(function(data, status) {
        var d = data.Data;
        extendSessionData(d);
        continueAfterSetUp(d);
    });
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
    var p = game.Players[myPlayer];
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
                confirm = $('#my-player .hand li').last();
                $.ajax({
                    url: '/3-13/test/' + myName + '/?move=' + CV(m),
                    dataType: 'json',
                }).done(function(data, status) {
                    if (data.Status !== 0) {
                        alert(data.Info);
                        cmdQ.push(synchronizeTableStatus);
                        return; // TODO reload?
                    }
                    confirm.html(cardToHtml(data.Data));
					confirm.find('.card').data('card', data.Data);
                    $('#tt-pile .card,#tt-deck .card').draggable('destroy');
                    cmdQ.push(synchronizeTableStatus);
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
                var c = d.find('.card').data('card');
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
                    cmdQ.push(synchronizeTableStatus);
                });
            },
        });
    } else if (p.State == CV('DONE')) {
        setPrompt('Set up groups or pass the turn');
        
        (function(){
            var a = $('#my-player .actions');
            a.show();
            a.find('.pass')
			.removeAttr('disabled')
			.click(function(ev) {
                $.ajax({
                    'url': '/3-13/test/' + myName +'/?move=' + CV('PASS_TURN'),
                    'dataType': 'json',
                }).done(function(data, status) {
                    if (data.Status !== 0) {
                        alert(data.Info);
                        return; // TODO reload?
                    }
                    cmdQ.push(synchronizeTableStatus);
                });
                ev.preventDefault();
				a.find('.pass').attr('disabled', 'disabled');
            });
            a.find('.done')
            .removeAttr('disabled')
            .click(function(ev) {
                var gs, hand;
                hand = [];
                $('.hand .card').each(function (idx, el) {
                    el = $(el);
                    hand.push(el.data('card'));
                });
                a.find('.pass,.done').attr('disabled','disabled');

                gs = $('<div class="group-setup">'+
                        '<div class="groups-wrap"><div class="groups"/></div>'+
                        '<input class="add-group" type="button" value="+">'+
                        '<input class="cancel-groups" type="button" value="Cancel">'+
                        '<input class="send-groups" type="button" value="Send">'+
                        '</div>');
                gs.find('.add-group').click(function(ev) {
                    addGroup();
                    ev.preventDefault();
                });
                gs.find('.cancel-groups').click(function(ev) {
                    var i,h;
                    $('.hand li').remove();
                    h = $('.hand');
                    for (i in hand) {
                        addCardItem(h, hand[i]);
                    }
                    $('.group-setup').remove();
                    a.find('.pass,.done').removeAttr('disabled');
                    ev.preventDefault();
                });
                gs.find('.send-groups').click(function(ev) {
                    var gStr, uStr;
                    gStr = '';
                    $('.groups .group').each(function(idx,el) {
                        $(el).find('.card').each(function(idx,el) {
                            gStr += cardToAscii($(el).data('card')) + ',';
                        });
                        gStr = gStr.substr(0, gStr.length - 1);
                        gStr += ';';
                    });
                    gStr = gStr.substr(0, gStr.length - 1);

                    console.log(gStr);

                    uStr = '';
                    $('.hand .card').each(function (idx, el) {
                        uStr += cardToAscii($(el).data('card')) + ',';
                    });
                    uStr = uStr.substr(0, uStr.length - 1);

                    $.ajax({
                        'url': '/3-13/test/' + myName + '/?move=' + CV('DECLARE_DONE') +
                            '&groups=' + gStr + '&unassigned=' + uStr,
                        'dataType': 'json',
                    }).done(function (data, status) {
                        if (data.Status !== 0) {
                            alert(data.Info);
                            return; // TODO reload?
                        }
                        $('.group-setup').remove();
                        cmdQ.push(synchronizeTableStatus);
                    });
                    ev.preventDefault();
                });
                $('#my-player').prepend(gs);
                addGroup();
            });
        })();
    } else {
        setPrompt(p.State);
    }
}

function addGroup() {
    var g;
    g = $('<div class="group"><input class="remove-group" type="button" value="X"><ul></ul></div>');
    g.find('ul').sortable({
        'connectWith': 'ul.hand, .group ul',
        'placeholder': 'hand-placeholder',
    });
    g.find('input.remove-group').click(function(ev) {
        var me = $(this).parent();
        me.find('ul li').each(function(idx, el) {
            $(el).detach().appendTo('ul.hand');
        });
        me.remove();
        ev.preventDefault();
    });
    $('.group-setup .groups').append(g);
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
        d.Game.Players[i].FinalGroup = d.Game.FinalGroups[i];
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
	var i, p;
	$('.active-player').removeClass('active-player');
	$(pId(game.CurrentPlayer)).addClass('active-player');

    for (i=0; i<game.Players.length; ++i) {
        if (myPlayer == i) continue;
        p = game.Players[i];
        updateCardsInHand(p);
        checkAndSetFinalGroups(p);
	}

	updateMyHand();
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

function checkAndSetFinalGroups(player) {
    var fg,p,i,f,set,unassigned;
    if (player.FinalGroup.Set === null) {
        return;
    }
    fg = assureFinalGroupsDivExists(pId(player.Index));

    p = $(pId(player.Index));
    f = p.find('ul.fan');
    f.html();
    addAllCardsTo(f, player.FinalGroup.Unassigned);

    fg.html();
    set = player.FinalGroup.Set;
    for (i in set) {
        f = $('<ul class="fan"/>');
        addAllCardsTo(f, set[i]);
        fg.append(f);
    }
    console.log(player);
}

function assureFinalGroupsDivExists(id) {
    var fg, fgSel;
    fgSel = id + ' .final-groups';
    fg = $(fgSel);
    if (fg.length === 0) {
        $(id).append('<div class="final-groups"/>');
        fg = $(fgSel);
    }
    return fg;
}

function addMyPlayer(hand) {
	var h, f, a, i;

    $('#my-player').append(
            playerHtml(myName),
            '<div class="actions"></div>',
            '<div class="prompt"></div>'
    );
	h = $('#my-player .player-box');
	h.attr('id', 'player-' + myPlayer);

    a = $('#my-player .actions');
    a.hide();
    a.append('<input class="deal" type="button" value="Deal"/>');
    a.append('<input class="done" type="button" value="Groups"/>');
    a.append('<input class="pass" type="button" value="Pass"/>');
    a.find('input').attr('disabled','disabled');

	f = h.find('ul.fan');
	f.removeClass('fan');
	f.addClass('hand');
    f.wrap('<div class="hwrap"/>');
	f.sortable({
    	placeholder: 'hand-placeholder',
        connectWith: '.group ul',
        revert: 250,
    });
    f.disableSelection();
 
	if (hand === null) {
		return;
	}
	for (i=0; i<hand.length;++i) {
	    addCardItem(f, hand[i]);
    }
}

function addCardItem(to, card) {
    to.append('<li class="dense">'+ cardToHtml(card) +'</li>');
    to.find('li').last().find('.card').data('card', card);
    return to;
}

function addAllCardsTo(to, arr) {
    var i;
    for (i in arr) {
        addCardItem(to, arr[i]);
    }
    return to;
}

function addMyHand() {
	return $.ajax({
		url: '/3-13/test/'+myName+'/?hand',
		dataType: 'json',
	}).done(function(data, status) {
		addMyPlayer(data.Data);
	});
}

function updateMyHand() {
	return $.ajax({
		url: '/3-13/test/'+myName+'/?hand',
		dataType: 'json',
	}).done(function(data, status) {
		var hand = data.Data, i, f;

		f = $('#my-player .hand');
		f.find('li').each(function(idx, el) {
			var i, hit, c;
			el = $(el);
            c = el.find('.card');
			hit = false;
			for (i=0; i < hand.length; ++i) {
				if (cardsEqual(c.data('card'), hand[i])) {
					hand.splice(i, 1);
					hit = true;
					break;
				}
			}
			if (!hit) {
				console.log('not hit ' + cardToAscii(c.data('card')));
				el.remove();
			}
		});

		for (i in hand) {
            addCardItem(f, hand[i]);
		}
	});
}

