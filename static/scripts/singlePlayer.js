const socket = new WebSocket("ws://localhost:2357/wsSinglePlayer");
const game_id = document.getElementById("game_id").innerText;

let opening = true;
let width,height,tot_bombs,time;

function setup() {
    const section = document.querySelector("#gameSection");
    socket.addEventListener("open", ()=>{
        socket.send(game_id);
        const canvas = createCanvas(200, 200);
        canvas.parent(section);
        resizeCollapsable();
        
        opening=false;
    });
}

function draw() {
    if (opening) return;
    background(59);
}

socket.addEventListener("message", e=>{
    e.data.arrayBuffer().then(ab=>{
        const data_view = new DataView(ab);
        for (let i=0; i<data_view.byteLength/2; i++) {
            console.log(data_view.getUint16(i*2));
        }
    });
});
