// var joy = new JoyStick('dpad', {"internalFillColor": "#DCDCDC", "externalStrokeColor": "black"});

const ws = new WebSocket("ws://" + document.domain + ':' + location.port + "/xbox");
// var joy = new JoyStick('Axis', {internalFillColor: "rgb(48, 47, 47)", internalStrokeColor: "black", externalStrokeColor: "black"})
var joy = new JoyStick('AxisR', {internalFillColor: "rgb(48, 47, 47)", internalStrokeColor: "black", externalStrokeColor: "black"})

var dpadVal = 0;

ws.addEventListener('open', function (event) {});

window.onload = function() {
    fetch('dpad.svg').then(response => response.text())
    .then((data) => {
        document.getElementById('Dpad').innerHTML += data
    }).then(() => setupbuttons());
}

function setupbuttons() {
    var anchors = document.getElementsByTagName('a');
    for (var i = 0; i < anchors.length; i++) {
        const anchor = anchors[i];
        const id = anchor.id;
        const cls = anchor.classList[0];

        if(cls=='Btn' || cls=='Trigger'){
            anchor.ontouchstart = function (ev) {
                ev.preventDefault();
                ev.stopPropagation();
                ws.send(`${cls}${id}|1`);   //eg: BtnA|1, TriggerR|1
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
    var paths = document.getElementsByTagName('path');
    console.log(paths);
    for (var i = 0; i < paths.length; i++) {
        var path = paths[i];
        const weight = path.dataset.weight;
        path.ontouchstart = function (ev) {
            ev.preventDefault();
            ev.stopPropagation();
            dpadVal += +weight;
            ws.send(`Dpad|${dpadVal}`);
            path.style.top = "3px";
        }
        path.ontouchend = function (ev) {
            ev.preventDefault();
            dpadVal -= +weight;
            ws.send(`Dpad|${dpadVal}`);
            path.style.top = "0px";
        }
    }
}

