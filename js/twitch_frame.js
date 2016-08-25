window.addEventListener('load', function () {
    'use strict';
    var tftv1 = 'https://api.twitch.tv/kraken/streams/teamfortresstv';
    var tftv2 = 'https://api.twitch.tv/kraken/streams/teamfortresstv2';

    var tftv1_player = 'https://player.twitch.tv/?channel=teamfortresstv';
    var tftv2_player = 'https://player.twitch.tv/?channel=teamfortresstv2';

    var cur_stream = 1;

    $.getJSON(tftv1, function(data) {
	if (data.stream == null) {
	    console.log("TFTV1 not live");

	    $.getJSON(tftv2, function(data) {
		if (data.stream != null) {
		    change_stream_url(tftv2_player, 1, 2);
		} else {
		    console.log("TFTV2 not live either :(");
		}
	    })
	}
    });

    var stream_frame = document.getElementById('stream');
    var button = $('#stream-button');

    button.click(function() {
	if (cur_stream == 1) {
	    change_stream_url(tftv2_player, 2, 1);
	} else {
	    change_stream_url(tftv1_player, 1, 2);
	}
    });

    function change_stream_url(url, cur_n, prev_n) {
	if (cur_stream == cur_n) {
	    return;
	}

	cur_stream = cur_n;
	console.log("Switching to "+url);
	stream_frame.src = url;
	button.text("teamfortresstv #"+prev_n);
    }
});
