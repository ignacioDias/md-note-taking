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

const $registerButton = document.querySelector(".register-button")
const $loginButton = document.querySelector(".login-button")

$registerButton.addEventListener("click", () => {
    window.location.href = '/register';
})
$loginButton.addEventListener("click", () => {
    window.location.href = '/login';
})