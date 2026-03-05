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


const $RegisterButton = document.querySelector(".register-button")
const $form = document.querySelector(".form-login")
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

$RegisterButton.addEventListener("click", () => {
    window.location.href = '/register';
})

$form.addEventListener("submit", async (event) => {
    event.preventDefault()
    clearError()

    const formData = new FormData($form)
    const data = {
        email: formData.get("email"),
        password: formData.get("password")
    }

    try {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })

        if (response.ok) {
            window.location.href = "/dashboard"
        } else {
            const errorText = await response.text()
            showError(errorText)
        }
    } catch (error) {
        showError("Error al conectar con el servidor. Verifica tu conexión a internet")
        console.error(error)
    }
})