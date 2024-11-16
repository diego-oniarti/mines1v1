document.getElementById('themeButton').addEventListener('click', e => {
    const body = document.getElementsByTagName('body')[0];
    const isDarkMode = body.classList.contains('dark');
    body.classList.toggle('dark');
    document.getElementById('themeButton').innerHTML = isDarkMode ? '◑' : '◐';
    const expiryDate = new Date();
    expiryDate.setDate(expiryDate.getDate() + 3000);
    document.cookie = `mode=${isDarkMode ? 'light' : 'dark'};path=/;expires=${expiryDate.toUTCString()};SameSite=Lax`;
});
document.getElementById('languageButton').addEventListener('click',e=>{
    alert("This page is not translated yet");
});
//window.addEventListener("load", () => {
document.addEventListener("DOMContentLoaded", () => {
    const cookie = document.cookie;
    mode = cookie.match('mode=(?<mode>light|dark)')?.groups['mode'];
    console.log(mode)
    if (mode==="light") {
        document.getElementsByTagName('body')[0].classList.remove('dark');
    }
    var sheet = document.styleSheets[0];
    var rules = sheet.cssRules || sheet.rules;

    //rules[0].style.color = 'red';
    setTimeout(()=>{
        rules[0].style.transition = 'color 200ms ease-in-out, background-color 200ms ease-in-out';
    },500);
});

document.getElementById('navCollapse').addEventListener('click',e=>{
    for (element of document.getElementsByClassName('navCollapsible')) {
        if (element.style.maxHeight){
            element.style.maxHeight = null;
        } else {
            element.style.maxHeight = element.scrollHeight + "px";
        }
    }
});

function resizeCollapsable () {
    for (element of document.getElementsByClassName('section')) {
        if (element.style.maxHeight!='0px')
            element.style.maxHeight = element.scrollHeight + "px";
    }
}
resizeCollapsable();
window.addEventListener('resize', ()=>{
    resizeCollapsable();
})

for (let bottone of document.getElementsByClassName('collapseButton')) {    
    bottone.addEventListener('click', ()=>{
        bottone.children[1].classList.toggle('flipped')
        for (element of bottone.parentElement.children) {
            if (element.style.backgroundColor=='transparent') 
                element.style.backgroundColor=''
            else
                element.style.backgroundColor='transparent'
            
            if (element.classList.contains('section')) {
                if (element.style.maxHeight!='0px'){
                    element.style.maxHeight = '0px';
                } else {
                    element.style.maxHeight = element.scrollHeight + "px";
                }
            }
        }
    });
}

for (let toggle of document.getElementsByClassName("toggle")) {
    toggle.addEventListener("click", ()=>{
        toggle.classList.toggle("toggled");
        for (let child of toggle.children) {
            if (child.tagName==="INPUT") {
                child.checked = !child.checked;
            }
        }
    })
}
