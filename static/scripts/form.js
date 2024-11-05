function validaLogin(form) {
    /** @type {string} */
    const mail = form.mail.value.trim();
    /** @type {string} */
    const psw = form.password.value.trim();

    const err_box = document.getElementById("login_error_box");
    const err_p = document.getElementById("login_error_p");

    const errors = [];

    if (!mail.match(/.+@.+\..+/)) {
        errors.push("Mail non valida")
    }

    if ((psw.length<8)
        || (!psw.match(/[!@#$%^&*()_+\-=\[\]{}]/))
        || (!psw.match(/[a-z]/))
        || (!psw.match(/[A-Z]/))
        || (!psw.match(/[0-9]/))) {
        errors.push("Password non valida")
    }

    if (errors.length>0) {
        err_box.classList.remove("hidden");
        err_p.innerHTML = errors.map(v=>"- "+v).join("<br/>");
        resizeCollapsable();
        return false;
    }

    return true;
}

function validaRegistrazione(form) {
    /** @type {string} */
    const nome = form.username.value.trim();
    /** @type {string} */
    const mail = form.mail.value.trim();
    /** @type {string} */
    const psw = form.password.value.trim();
    /** @type {string} */
    const psw2 = form.password_confirm.value.trim();

    const err_box = document.getElementById("register_error_box");
    const err_p = document.getElementById("register_error_p");

    const errors = [];

    if (nome.length<1) {
        errors.push("L'username deve avere almeno 1 carattere");
    }
    if (!mail.match(/.+@.+\..+/)) {
        errors.push("Inserisci una mail valida");
    }
    if (psw.length<8) {
        errors.push("La passwoed deve contenere almeno 8 caratteri");
    }
    if (!psw.match(/[_!?(){}#$%^&*.,+\[\]=+"']/)) {
        console.log(psw)
        errors.push("La password deve contenere almeno 1 carattere speciale (!@#$%^&*()_+-=[]{})");
    }
    if (!psw.match(/[a-z]/)) {
        errors.push("La password deve contenere almeno 1 carattere minuscolo");
    }
    if (!psw.match(/[A-Z]/)) {
        errors.push("La password deve contenere almeno 1 carattere maiuscolo");
    }
    if (!psw.match(/[0-9]/)) {
        errors.push("La password deve contenere almeno 1 numero");
    }
    if (psw!=psw2) {
        errors.push("Conferma password");
    }


    if (errors.length>0) {
        err_box.classList.remove("hidden");
        err_p.innerHTML = errors.map(v=>"- "+v).join("<br/>");
        resizeCollapsable();
        return false;
    }

    return true;
}
