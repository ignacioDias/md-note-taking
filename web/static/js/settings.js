// Aplicar el color de fondo guardado al cargar la página
function applyBackgroundColor() {
    const savedColor = localStorage.getItem('backgroundColor');
    if (savedColor) {
        document.body.style.backgroundColor = savedColor;
        
        // Marcar el preset activo si coincide
        document.querySelectorAll('.color-preset').forEach(btn => {
            if (btn.dataset.color === savedColor) {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });
        
        // Actualizar el color personalizado
        document.getElementById('customColor').value = savedColor;
    }
}

// Guardar y aplicar color de fondo
function setBackgroundColor(color) {
    localStorage.setItem('backgroundColor', color);
    document.body.style.backgroundColor = color;
}

// Inicializar al cargar la página
document.addEventListener('DOMContentLoaded', () => {
    applyBackgroundColor();
    
    // Manejar clics en los presets de color
    document.querySelectorAll('.color-preset').forEach(btn => {
        btn.addEventListener('click', () => {
            const color = btn.dataset.color;
            setBackgroundColor(color);
            
            // Marcar como activo
            document.querySelectorAll('.color-preset').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            // Actualizar el selector de color personalizado
            document.getElementById('customColor').value = color;
        });
    });
    
    // Aplicar color personalizado
    document.getElementById('applyCustomColor').addEventListener('click', () => {
        const color = document.getElementById('customColor').value;
        setBackgroundColor(color);
        
        // Quitar la marca activa de los presets
        document.querySelectorAll('.color-preset').forEach(b => b.classList.remove('active'));
    });
    
    // Restablecer a color predeterminado
    document.getElementById('resetColor').addEventListener('click', () => {
        localStorage.removeItem('backgroundColor');
        document.body.style.backgroundColor = '#ffffff';
        document.getElementById('customColor').value = '#ffffff';
        
        // Marcar el preset blanco como activo
        document.querySelectorAll('.color-preset').forEach(btn => {
            if (btn.dataset.color === '#ffffff') {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });
    });
    
    // Navegación
    document.querySelector('.dashboard-btn').addEventListener('click', () => {
        window.location.href = '/dashboard';
    });
    
    document.querySelector('.logout-btn').addEventListener('click', async () => {
        try {
            const response = await fetch('/logout', { method: 'POST' });
            if (response.ok) {
                window.location.href = '/';
            }
        } catch (error) {
            console.error('Error al cerrar sesión:', error);
        }
    });
});
