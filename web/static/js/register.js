document.addEventListener('DOMContentLoaded', async function() {
    try {
        const response = await fetch('/api/me', {
            method: 'GET',
            credentials: 'include'
        });
        if (response.ok) {
            window.location.href = '/dashboard';
        }
    } catch (error) {
        console.error('Error checking authentication:', error);
    }
});


const $LoginButton = document.querySelector(".login-button")
const $form = document.querySelector(".form-register")
const $errorContainer = document.getElementById("error-container")

function showError(message) {
    $errorContainer.innerHTML = `
        <div class="alert alert-danger alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>
    `
}

function clearError() {
    $errorContainer.innerHTML = ""
}

$LoginButton.addEventListener("click", () => {
    window.location.href = '/login';
})

$form.addEventListener("submit", async (event) => {
    event.preventDefault()
    clearError()
    const formData = new FormData($form)
    if (formData.get("password") != formData.get("repeated-password")) {
        showError("Contraseñas distintas")
        return
    }
    const data = {
        name: formData.get("name"),
        email: formData.get("email"),
        password: formData.get("password")
    }

    try {
        const response = await fetch("/api/auth/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })

        if (response.ok) {
            window.location.href = "/login"
        } else {
            const errorText = await response.text()
            showError(errorText)
        }
    } catch (error) {
        showError("Error al conectar con el servidor. Verifica tu conexión a internet")
        console.error(error)
    }
})