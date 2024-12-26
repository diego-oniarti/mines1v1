const socket = new WebSocket("ws://localhost:2357/wsSinglePlayer");
const game_id = document.getElementById("game_id").innerText;
const section = document.querySelector("#gameSection");
const box = document.querySelector("#gameBox");

const endgame_controls = document.querySelector("#endgame_controls");
const play_again = document.querySelector("#replay");
const goback = document.querySelector("#goback");

socket.addEventListener("open", ()=>{
    socket.send(game_id);
});

let grid_width,grid_height,tot_bombs,time;
let a; let b = new Promise(r=>{a=r});
let cellSize=40;

class Cella {
    constructor(flag, number, bomb) {
        this.flag = flag;
        this.number = number;
        this.bomb = bomb;
    }
}
/** @type{Cella[][]} */
const celle = [];

async function setup() {
    await b;
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
    canvas.parent(box);
    section.addEventListener('contextmenu', event => event.preventDefault());
    resizeCollapsable();
}

let timer_start;
function draw() {
    if (phase==phases.GetGameParams) return;
    background(200);
    stroke(150);
    strokeWeight(1);
    for (let i=1; i<grid_width; i++) {
        line(i*cellSize,0, i*cellSize,height);
    }
    for (let i=1; i<grid_height; i++) {
        line(0,i*cellSize, width, i*cellSize);
    }

    textAlign(CENTER, CENTER);
    textSize(cellSize*0.8);
    noStroke();
    for (let y=0; y<grid_height; y++) {
        for (let x=0; x<grid_width; x++) {
            if (!celle[y][x]) continue;
            if (celle[y][x].flag) {
                text("🏳", (x+0.5)*cellSize, (y+0.5)*cellSize);
                continue;
            }
            fill(150);
            rect(x*cellSize, y*cellSize, cellSize, cellSize);
            fill(0);
            if (celle[y][x].bomb) {
                stroke(200,0,0);
                line(x*cellSize, y*cellSize,(x+1)*cellSize, (y+1)*cellSize);
                line((x+1)*cellSize, y*cellSize,x*cellSize, (y+1)*cellSize);
                noStroke();
                continue;
            }
            if (celle[y][x].number!=0) {
                text(celle[y][x].number, (x+0.5)*cellSize, (y+0.5)*cellSize);
            }
        }
    }

    switch (phase) {
        case phases.GetUpdates:
            if (time==0 || !timer_start) break;
            const R = Math.min(width,height)*0.8
            const c = color(180,40,40,30);
            stroke(c);
            noFill();
            strokeWeight(20);
            circle(width/2,height/2, R);
            noStroke();
            fill(c);
            arc(width/2, height/2, R-20, R-20, map(new Date()-timer_start, 0, time, 0,TWO_PI)-PI/2, -PI/2);
            break;
        case phases.Won:
            textSize(height/5);
            textStyle(BOLD);
            fill(50,200,50);
            stroke(50,100,50);
            strokeWeight(2);
            text("YOU WON", width/2, height/2);
            textStyle(NORMAL);
            break;
        case phases.Lost:
            textSize(height/5);
            textStyle(BOLD);
            fill(200,0,0);
            stroke(100,50,50);
            strokeWeight(2);
            text("YOU LOST", width/2, height/2);
            textStyle(NORMAL);
            break;
    }
}

function mousePressed() {
    if (mouseX>width||mouseX<0||mouseY<0||mouseY>height) return true;
    if (phase != phases.GetUpdates) return true;
    const x = Math.floor(mouseX/cellSize);
    const y = Math.floor(mouseY/cellSize);
    const flag = mouseButton!=LEFT;
    const bits = new ArrayBuffer(5);
    const view = new DataView(bits);
    view.setInt16(0, x);
    view.setInt16(2, y);
    view.setInt8(4, flag?1:0);

    socket.send(bits);

    return false;
}

const phases = {
    GetGameParams: 1,
    GetUpdates: 2,
    Won: 3,
    Lost: 4,
}
let phase = phases.GetGameParams;

socket.addEventListener("message", e=>{
    e.data.arrayBuffer().then(ab=>{
        const data_view = new DataView(ab);

        switch (phase) {
            case phases.GetGameParams:
                get_game_params(data_view);
                break;
            case phases.GetUpdates:
                get_updates(data_view);
                break;
        }
    });
});

function get_game_params(data_view) {
    const data = [];
    for (let i=0; i<4; i++) {
        data.push(data_view.getUint16(i*2))
    }
    [grid_width,grid_height,tot_bombs,time] = data;
    a();
    phase = phases.GetUpdates;
    for (let y=0; y<grid_width; y++) {
        celle.push([])
        for (let x=0; x<grid_width; x++) {
            celle[y].push(null)
        }
    }
}

/**
 * @param {DataView} data_view 
 */

function get_updates(data_view) {
    const first_byte = data_view.getInt8(0);
    type = first_byte >> 6;
    switch (type) {
        case 0:
            const gameover = (first_byte & 0b00100000)>0;
            const won = (first_byte & 0b00010000)>0;
            const lost = gameover && !won;
            let off = 1;
            let has_next = true;
            do {
                const x = data_view.getUint16(off); off+=2;
                const y = data_view.getUint16(off); off+=2;
                const payload = data_view.getUint8(off); off+=1;
                const num = payload>>4;
                has_next = (payload&0b00001000)>0;

                if (lost) {
                    celle[y][x] = new Cella(false, 0, true);
                }else{
                    celle[y][x] = new Cella(false, num, false);
                }
            } while (has_next);

            if (gameover) {
                phase = won ? phases.Won : phases.Lost;
                show_endgame_controls();
                return;
            } else {
                timer_start = new Date();
            }
            break;
        case 1:
            const flag = (first_byte&0b00100000)>0;
            const x = data_view.getUint16(1);
            const y = data_view.getUint16(3);
            if (flag) {
                celle[y][x] = new Cella(true, 0, false);
            }else{
                celle[y][x] = null;
            }
            break;
    }
}

function show_endgame_controls() {
    endgame_controls.classList.add("form_fields");
    resizeCollapsable();
}

play_again.addEventListener("click", ()=>{
    phase = phases.GetUpdates;
    // location.reload();
    for (let row of celle) {
        for (let i=0; i<row.length; i++) {
            row[i] = null;
        }
    }
    const a = new ArrayBuffer(1);
    (new DataView(a)).setUint8(0,1);
    socket.send(a);
});
