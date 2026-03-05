document.addEventListener('DOMContentLoaded', async function() {
    const response = await fetch('/api/me', { credentials: 'include' });
    if (!response.ok) {
        window.location.href = '/login';
    }
});