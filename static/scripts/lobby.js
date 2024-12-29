async function get_game_id(form) {
    const width = form.width.value;
    const height = form.height.value;
    const bombs = form.bombs.value;
    const tempo = form.tempo.value;
    const timed = form.timed.checked;

    const body = JSON.stringify({
        "width": parseInt(width),
        "height": parseInt(height),
        "bombs": parseInt(bombs),
        "tempo": parseInt(tempo),
        "timed": timed?"on":"off",
    });

    return await fetch("/createGame", {
        method: "POST",
        body: body,
    })
        .then(r=>r.json())
        .then(e=>{
            return e.game_id;
        });
}

function submitSingle(form) {
    get_game_id(form).then(id=>{
        window.location.href = `/singlePlayer?game_id=${id}`;
    });
    return false;
}
