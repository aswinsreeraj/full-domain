async function handleLogout() {
    const response = await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include'
    });

    if (response.redirected) {
        window.location.href = response.url;
        return;
    }

    window.location.href = "/";
}

async function handleAdminLogout() {
    const response = await fetch('/admin/logout', {
        method: 'POST',
        credentials: 'include'
    });

    if (response.redirected) {
        window.location.href = response.url;
        return;
    }

    window.location.href = "/";
}

function validatePasswordFields(passwordId, confirmId, required = false) {
    const password = document.getElementById(passwordId).value;
    const confirm = document.getElementById(confirmId).value;

    const regex =
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*(),.?":{}|<>_\-])[A-Za-z\d!@#$%^&*(),.?":{}|<>_\-]{8,}$/;

    if (!required && password.length === 0 && confirm.length === 0) {
        return true;
    }
    if (!regex.test(password)) {
        alert("Password must be at least 8 chars and include uppercase, lowercase, digit, and special char.");
        return false;
    }

    if (password !== confirm) {
        alert("Passwords do not match.");
        return false;
    }

    return true;
}


document.getElementById("signup-form").addEventListener("submit", function (e) {
    if (!validatePasswordFields("password", "confirm-password", true)) {
        e.preventDefault();
    }
});

document.getElementById("update-user").addEventListener("submit", function (e) {
    if (!validatePasswordFields("new-password", "confirm-password", false)) {
        e.preventDefault();
    }
});

document.querySelector('form[action="/admin/create"]').addEventListener("submit", function (e) {
    if (!validatePasswordFields("new_user_password", "confirm_new_user_password", true)) {
        e.preventDefault();
    }
});

function toggleCreate() {
    const url = new URL(window.location.href);

    if (url.searchParams.get("create") === "true") {
        url.searchParams.set("create", "false");
    } else {
        url.searchParams.set("create", "true");
    }

    window.location.href = url.toString();
}