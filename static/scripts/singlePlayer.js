const socket = new WebSocket("ws://localhost:2357/wsSinglePlayer");
const game_id = document.getElementById("game_id").innerText;
const section = document.querySelector("#gameSection");
socket.addEventListener("open", ()=>{
    socket.send(game_id);
});

let opening = true;
let grid_width,grid_height,tot_bombs,time;
let a; let b = new Promise(r=>{a=r});
let cellSize=30;

async function setup() {
    await b;
    // const AR = width/height;
    // let W = section.clientWidth;
    // let H = section.clientWidth/AR;
    // if (H>window.innerHeight*0.8) {
    //     H = window.innerHeight*0.8;
    //     W = H*AR;
    // }
    let W = grid_width*cellSize;
    let H = grid_height*cellSize;
    if (W > section.clientWidth) {
        cellSize = section.clientWidth/grid_width;
        W = grid_width*cellSize;
        H = grid_height*cellSize;
    }
    if (H>window.innerHeight*0.9) {
        cellSize = window.innerHeight*0.9/grid_height;
        W = grid_width*cellSize;
        H = grid_height*cellSize;
    }

    const canvas = createCanvas(W, H);
    canvas.parent(section);
    resizeCollapsable();
}

function draw() {
    if (opening) return;
    background(59);
    stroke(150);
    for (let i=0; i<grid_width; i++) {
        line(i*cellSize,0, i*cellSize,height);
    }
    for (let i=0; i<grid_height; i++) {
        line(0,i*cellSize, width, i*cellSize);
    }
}

function mouseClicked() {
    console.log("A");
    const x = Math.floor(mouseX/cellSize);
    const y = Math.floor(mouseY/cellSize);
    socket.send(new Uint16Array([x,y]).buffer);
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
                [grid_width,grid_height,tot_bombs,time] = data;
                a();
                opening = false;
                break;
        }
    });
});
