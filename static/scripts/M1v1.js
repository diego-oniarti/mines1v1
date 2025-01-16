const socket = new WebSocket("/ws1v1");
const game_id = document.getElementById("game_id").innerText;
const section = document.querySelector("#gameSection");
const box = document.querySelector("#gameBox");

const endgame_controls = document.querySelector("#endgame_controls");
const play_again = document.querySelector("#replay");
const goback = document.querySelector("#goback");

socket.addEventListener("open", () => {
    socket.send(game_id);
});

let grid_width, grid_height, tot_bombs, time;
let turn = false;

let a;
let b = new Promise((r) => {
    a = r;
});
let cellSize = 40;

let placed_bombs = 0;
const bomb_span = document.getElementById("bomb_span");
function update_bombs() {
    bomb_span.innerText = `🏳${placed_bombs}/${tot_bombs}`;
}
function sub_flag() {
    placed_bombs--;
    update_bombs();
}
function add_flag() {
    placed_bombs++;
    update_bombs();
}

class Cella {
    constructor(flag, number, bomb, player) {
        this.flag = flag;
        this.number = number;
        this.bomb = bomb;
        this.player = player;
    }
}
/** @type{Cella[][]} */
const celle = [];

let circles = [];

async function setup() {
    await b;
    let W = grid_width * cellSize;
    let H = grid_height * cellSize;
    if (W > section.clientWidth) {
        cellSize = section.clientWidth / grid_width;
        W = grid_width * cellSize;
        H = grid_height * cellSize;
    }
    if (H > window.innerHeight * 0.9) {
        cellSize = (window.innerHeight * 0.9) / grid_height;
        W = grid_width * cellSize;
        H = grid_height * cellSize;
    }

    const canvas = createCanvas(W, H);
    canvas.parent(box);
    section.addEventListener("contextmenu", (event) => event.preventDefault());
    resizeCollapsable();
}

let timer_start;
function draw() {
    if (phase == phases.GetGameParams) return;
    background(240);
    stroke(0);
    strokeWeight(0.1);
    for (let i = 1; i < grid_width; i++) {
        line(i * cellSize, 0, i * cellSize, height);
    }
    for (let i = 1; i < grid_height; i++) {
        line(0, i * cellSize, width, i * cellSize);
    }

    textAlign(CENTER, CENTER);
    textSize(cellSize * 0.8);
    noStroke();
    for (let y = 0; y < grid_height; y++) {
        for (let x = 0; x < grid_width; x++) {
            if (!celle[y][x]) continue;
            const cella = celle[y][x];
            noStroke();
            if (cella.flag) {
                if (cella.player) {
                    fill(0, 255, 255);
                } else {
                    fill(255, 0, 255);
                }
                text("🏳", (x + 0.5) * cellSize, (y + 0.5) * cellSize);
                continue;
            }
            fill(220);
            rect(x * cellSize, y * cellSize, cellSize, cellSize);
            if (cella.player) {
                fill(0, 190, 190);
            } else {
                fill(190, 0, 190);
            }
            if (cella.bomb) {
                stroke(200, 0, 0);
                line(
                    x * cellSize,
                    y * cellSize,
                    (x + 1) * cellSize,
                    (y + 1) * cellSize,
                );
                line(
                    (x + 1) * cellSize,
                    y * cellSize,
                    x * cellSize,
                    (y + 1) * cellSize,
                );
                noStroke();
                continue;
            }
            if (cella.number != 0) {
                text(cella.number, (x + 0.5) * cellSize, (y + 0.5) * cellSize);
            }
        }
    }

    const now = new Date();
    const circle_max = 200;
    circles = circles.filter((c) => {
        return now - c.start < circle_max;
    });
    noFill();
    strokeWeight(4);
    for (let c of circles) {
        stroke(190, 40, 40, map(now - c.start, 0, circle_max, 255, 0));
        circle(
            (c.x + 0.5) * cellSize,
            (c.y + 0.5) * cellSize,
            map(now - c.start, 0, circle_max, 0, cellSize),
        );
    }

    switch (phase) {
        case phases.GetUpdates:
            if (time == 0 || !timer_start) break;
            const R = Math.min(width, height) * 0.8;
            const c = !turn ? color(180, 40, 40, 30) : color(40, 40, 180, 30);
            stroke(c);
            noFill();
            strokeWeight(20);
            circle(width / 2, height / 2, R);
            noStroke();
            fill(c);
            arc(
                width / 2,
                height / 2,
                R - 20,
                R - 20,
                map(new Date() - timer_start, 0, time, 0, TWO_PI) - PI / 2,
                -PI / 2,
            );
            break;
        case phases.Won:
            textSize(height / 5);
            textStyle(BOLD);
            fill(50, 200, 50);
            stroke(50, 100, 50);
            strokeWeight(2);
            text("YOU WON", width / 2, height / 2);
            textStyle(NORMAL);
            break;
        case phases.Draw:
            textSize(height / 5);
            textStyle(BOLD);
            fill(0, 0, 200);
            stroke(100, 50, 50);
            strokeWeight(2);
            text("DRAW", width / 2, height / 2);
            textStyle(NORMAL);
            break;
        case phases.Lost:
            textSize(height / 5);
            textStyle(BOLD);
            fill(200, 0, 0);
            stroke(100, 50, 50);
            strokeWeight(2);
            text("YOU LOST", width / 2, height / 2);
            textStyle(NORMAL);
            break;
    }
}
function mousePressed() {
    if (mouseX > width || mouseX < 0 || mouseY < 0 || mouseY > height)
        return true;
    if (phase != phases.GetUpdates) return true;
    const x = Math.floor(mouseX / cellSize);
    const y = Math.floor(mouseY / cellSize);
    const flag = mouseButton != LEFT;
    const bits = new ArrayBuffer(5);
    const view = new DataView(bits);
    view.setInt16(0, x);
    view.setInt16(2, y);
    view.setInt8(4, flag ? 1 : 0);

    socket.send(bits);

    return false;
}

const phases = {
    GetGameParams: 1,
    GetUpdates: 2,
    Won: 3,
    Lost: 4,
    Draw: 5,
};
let phase = phases.GetGameParams;

socket.addEventListener("message", (e) => {
    e.data.arrayBuffer().then((ab) => {
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
    for (let i = 0; i < 5; i++) {
        data.push(data_view.getUint16(i * 2));
    }
    let player_num;
    [grid_width, grid_height, tot_bombs, time, player_num] = data;
    if (player_num == 0) {
        turn = true;
    }
    a();
    phase = phases.GetUpdates;
    for (let y = 0; y < grid_width; y++) {
        celle.push([]);
        for (let x = 0; x < grid_width; x++) {
            celle[y].push(null);
        }
    }
    update_bombs();
}

/**
 * @param {DataView} data_view
 */

function get_updates(data_view) {
    const first_byte = data_view.getInt8(0);
    const type = first_byte >> 6;
    const player = (first_byte & 1) == 0;
    let first_x, first_y;
    switch (type) {
        case 0:
            const gameover = (first_byte & 0b00100000) > 0;
            const won = (first_byte & 0b00010000) > 0;
            const lost = gameover && !won;
            let off = 1;
            let has_next = true;

            let is_first = true;

            do {
                const x = data_view.getUint16(off);
                off += 2;
                const y = data_view.getUint16(off);
                off += 2;
                const payload = data_view.getUint8(off);
                off += 1;
                const num = payload >> 4;
                has_next = (payload & 0b00001000) > 0;

                if (lost) {
                    celle[y][x] = new Cella(false, 0, true, player);
                } else {
                    if (celle[y][x]?.flag) {
                        sub_flag();
                    }
                    celle[y][x] = new Cella(false, num, false, player);
                }

                if (is_first) {
                    first_x = x;
                    first_y = y;
                    is_first = false;
                }
            } while (has_next);

            if (gameover) {
                if (won) {
                    phase = phases.Draw;
                } else {
                    if (turn) {
                        phase = phases.Lost;
                    } else {
                        phase = phases.Won;
                    }
                }
                console.log("timer null");
                timer_start = null;
                show_endgame_controls();
                return;
            } else {
                timer_start = new Date();
            }

            turn = !turn;
            break;
        case 1:
            const flag = (first_byte & 0b00100000) > 0;
            const x = data_view.getUint16(1);
            const y = data_view.getUint16(3);
            first_x = x;
            first_y = y;
            if (flag) {
                celle[y][x] = new Cella(true, 0, false, player);
                add_flag();
            } else {
                celle[y][x] = null;
                sub_flag();
            }
            break;
    }

    if (!player) {
        circles.push({ x: first_x, y: first_y, start: new Date() });
    }
}

function show_endgame_controls() {
    endgame_controls.classList.add("form_fields");
    resizeCollapsable();
}

play_again.addEventListener("click", () => {
    for (let row of celle) {
        for (let i = 0; i < row.length; i++) {
            row[i] = null;
        }
    }

    placed_bombs = 0;
    update_bombs();
    socket.send("replay");
    phase = phases.GetUpdates;
});
