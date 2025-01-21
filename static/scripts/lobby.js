async function get_game_id(form) {
    const width = form.width.value;
    const height = form.height.value;
    const bombs = form.bombs.value;
    const tempo = form.tempo.value;
    const timed = form.timed.checked;

    const body = JSON.stringify({
        width: parseInt(width),
        height: parseInt(height),
        bombs: parseInt(bombs),
        tempo: parseInt(tempo),
        timed: timed ? "on" : "off",
    });

    return await fetch("/createGame", {
        method: "POST",
        body: body,
    })
        .then((r) => r.json())
        .then((e) => {
            return e.game_id;
        });
}

function submitSingle(form) {
    get_game_id(form).then((id) => {
        window.location.href = `/singlePlayer?game_id=${id}`;
    });
    return false;
}

function submit1v1(form) {
    get_game_id(form).then((id) => {
        window.location.href = `/1v1?game_id=${id}`;
    });
    return false;
}

document.addEventListener("DOMContentLoaded", () => {
    const form_single = document.querySelector("#form_single");
    const form_1v1 = document.querySelector("#form_1v1");
    document
        .querySelector("#single_difficulty")
        .addEventListener("change", (e) => {
            const width = form_single.width;
            const height = form_single.height;
            const bombs = form_single.bombs;
            const timed = form_single.timed;
            const tempo = form_single.tempo;
            switch (e.target.value) {
                case "easy":
                    width.value = 10;
                    height.value = 8;
                    bombs.value = 10;
                    timed.checked = true;
                    tempo.value = 10000;
                    break;
                case "medium":
                    width.value = 18;
                    height.value = 14;
                    bombs.value = 40;
                    timed.checked = true;
                    tempo.value = 3000;
                    break;
                case "hard":
                    width.value = 24;
                    height.value = 20;
                    bombs.value = 99;
                    timed.checked = true;
                    tempo.value = 1000;
                    break;
            }
            if (e.target.value == "custom") {
                width.disabled = false;
                height.disabled = false;
                bombs.disabled = false;
                timed.disabled = false;
                tempo.disabled = false;
            } else {
                width.disabled = true;
                height.disabled = true;
                bombs.disabled = true;
                timed.disabled = true;
                tempo.disabled = true;
            }
        });

    document
        .querySelector("#M1v1_difficulty")
        .addEventListener("change", (e) => {
            const width = form_1v1.width;
            const height = form_1v1.height;
            const bombs = form_1v1.bombs;
            const timed = form_1v1.timed;
            const tempo = form_1v1.tempo;
            switch (e.target.value) {
                case "easy":
                    width.value = 10;
                    height.value = 8;
                    bombs.value = 10;
                    timed.checked = true;
                    tempo.value = 10000;
                    break;
                case "medium":
                    width.value = 18;
                    height.value = 14;
                    bombs.value = 40;
                    timed.checked = true;
                    tempo.value = 3000;
                    break;
                case "hard":
                    width.value = 24;
                    height.value = 20;
                    bombs.value = 99;
                    timed.checked = true;
                    tempo.value = 1000;
                    break;
            }
            if (e.target.value == "custom") {
                width.disabled = false;
                height.disabled = false;
                bombs.disabled = false;
                timed.disabled = false;
                tempo.disabled = false;
            } else {
                width.disabled = true;
                height.disabled = true;
                bombs.disabled = true;
                timed.disabled = true;
                tempo.disabled = true;
            }
        });
});
