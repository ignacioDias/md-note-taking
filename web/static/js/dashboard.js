document.addEventListener('DOMContentLoaded', async function() {
    const response = await fetch('/api/me', { credentials: 'include' });
    if (!response.ok) {
        window.location.href = '/login';
    }
});

const $logoutBtn = document.querySelector(".logout-btn")
const $settingsBtn = document.querySelector(".settings-btn")
const $profileBtn = document.querySelector(".profile-btn")

$logoutBtn.addEventListener("click", async () => {
    try {
        const response = await fetch("/api/auth/logout", {
            method: "DELETE",
        })

        if (response.ok) {
            window.location.href = "/"
        } else {
            const errorText = await response.text()
            showError(errorText)
        }
    } catch (error) {
        showError("Error al conectar con el servidor. Verifica tu conexión a internet")
        console.error(error)
    }
})

$settingsBtn.addEventListener("click", () => {
    window.location.href = "/settings"
})
$profileBtn.addEventListener("click", () => {
    window.location.href = "/me"
})
