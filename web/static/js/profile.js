async function loadUserProfile() {
    try {
        const response = await fetch('/api/me', { credentials: 'include' });
        
        if (!response.ok) {
            window.location.href = '/login';
            return;
        }
        
        const userData = await response.json();
        
        document.getElementById('userName').textContent = userData.name || 'No disponible';
        document.getElementById('userEmail').textContent = userData.email || 'No disponible';
        
        if (userData.profilePicture || userData.profile_picture) {
            const profilePicUrl = userData.profilePicture || userData.profile_picture;
            const imgElement = document.getElementById('profilePicture');
            imgElement.src = profilePicUrl;
            imgElement.style.display = 'block';
        } else {
            document.getElementById('profilePicture').style.display = 'none';
        }
        
    } catch (error) {
        console.error('Error al cargar el perfil:', error);
        document.getElementById('userName').textContent = 'Error al cargar';
        document.getElementById('userEmail').textContent = 'Error al cargar';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    loadUserProfile();
    
    document.querySelector('.dashboard-btn').addEventListener('click', () => {
        window.location.href = '/dashboard';
    });
    
    document.querySelector('.settings-btn').addEventListener('click', () => {
        window.location.href = '/settings';
    });
    
    document.querySelector('.profile-btn').addEventListener('click', () => {
        window.location.reload();
    });
    
    document.querySelector('.logout-btn').addEventListener('click', async () => {
        try {
            const response = await fetch('/api/auth/logout', {
                method: 'DELETE',
                credentials: 'include'
            });
            
            if (response.ok) {
                window.location.href = '/';
            } else {
                console.error('Error al cerrar sesión');
            }
        } catch (error) {
            console.error('Error al cerrar sesión:', error);
        }
    });
    
    document.getElementById('profileSettingsBtn').addEventListener('click', () => {
        window.location.href = '/profile/settings';
    });
});
