fetch("https://horse-smooth-mutt.ngrok-free.app/minesredirect").then(res=>res.json()).then(data=>{
    console.log(data);
    window.location.replace(data.url);
})
