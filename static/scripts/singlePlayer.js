const socket = new WebSocket("ws://localhost:2357/wsSinglePlayer");
const game_id = document.getElementById("game_id").innerText;
const section = document.querySelector("#gameSection");
socket.addEventListener("open", ()=>{
    socket.send(game_id);
});

let opening = true;
let width,height,tot_bombs,time;
let a; let b = new Promise(r=>{a=r});

async function setup() {
    await b;
    const AR = width/height;
    let W = section.clientWidth;
    let H = section.clientWidth/AR;
    if (H>window.innerHeight*0.8) {
        H = window.innerHeight*0.8;
        W = H*AR;
    }

    const canvas = createCanvas(W, H);
    canvas.parent(section);
    resizeCollapsable();
}

function draw() {
    if (opening) return;
    background(59);
}

const phases = {
    GetGameParams: 1,
}
let phase = phases.GetGameParams;

socket.addEventListener("message", e=>{
    e.data.arrayBuffer().then(ab=>{
        const data_view = new DataView(ab);
        const data = [];
        for (let i=0; i<data_view.byteLength/2; i++) {
            console.log(data_view.getUint16(i*2));
            data.push(data_view.getUint16(i*2))
        }

        switch (phase) {
            case phases.GetGameParams:
                [width,height,tot_bombs,time] = data;
                a();
                opening = false;
                break;
        }
    });
});
