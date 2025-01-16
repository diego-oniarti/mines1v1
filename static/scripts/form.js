function validaLogin(form) {
    /** @type {string} */
    const mail = form.mail.value.trim();
    /** @type {string} */
    const psw = form.password.value.trim();

    const err_box = document.getElementById("login_error_box");
    const err_p = document.getElementById("login_error_p");

    const errors = [];

    if (!mail.match(/.+@.+\..+/)) {
        errors.push("Invalid email")
    }

    if ((psw.length<8)
        || (!psw.match(/[!@#$%^&*()_+\-=\[\]{}]/))
        || (!psw.match(/[a-z]/))
        || (!psw.match(/[A-Z]/))
        || (!psw.match(/[0-9]/))) {
        errors.push("Invalid password")
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
        errors.push("Username can't be empty");
    }
    if (!mail.match(/.+@.+\..+/)) {
        errors.push("Invalid email");
    }
    if (psw.length<8) {
        errors.push("Lassword must be at least 8 characters");
    }
    if (!psw.match(/[_!?(){}#$%^&*.,+\[\]=+"']/)) {
        console.log(psw)
        errors.push("Password must contain at least 1 special character (!@#$%^&*()_+-=[]{})");
    }
    if (!psw.match(/[a-z]/)) {
        errors.push("Password must contain at least 1 uppercase character");
    }
    if (!psw.match(/[A-Z]/)) {
        errors.push("Password must contain at least 1 lowercase character");
    }
    if (!psw.match(/[0-9]/)) {
        errors.push("Password must contain at least 1 number");
    }
    if (psw!=psw2) {
        errors.push("Confirm password");
    }


    if (errors.length>0) {
        err_box.classList.remove("hidden");
        err_p.innerHTML = errors.map(v=>"- "+v).join("<br/>");
        resizeCollapsable();
        return false;
    }

    return true;
}
