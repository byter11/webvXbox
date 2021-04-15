// var socket = io.connect(location.protocol + '//' + document.domain + ':' + location.port + '/socket.io/');
// var joy = new JoyStick('dpad', {"internalFillColor": "#DCDCDC", "externalStrokeColor": "black"});

const ws = new WebSocket("ws://" + document.domain + ':' + location.port + "/ter");
var joy = new JoyStick('Axis')

ws.addEventListener('open', function (event) {
});

window.onload = function() {
    setupbuttons();
}

function setupbuttons() {
    var anchors = document.getElementsByTagName('a');
    for (var i = 0; i < anchors.length; i++) {
        const anchor = anchors[i];
        const id = anchor.id;
        console.log(anchor);
        const cls = anchor.classList[0];

        if(cls=='Btn' || cls=='Trigger'){
            anchor.ontouchstart = function (ev) {
                ev.preventDefault();
                ev.stopPropagation();
                ws.send(`${cls}${id}|1`);
                anchor.style.position = "relative";
                anchor.style.top = "3px";
                
            }
            anchor.ontouchend = function (ev) {
                ev.preventDefault();
                ws.send(`${cls}${id}|0`);
                anchor.offsetTop = 0;
                anchor.style.top = "0px";
            }
        }
    }
}