function applyBackgroundColor() {
    const savedColor = localStorage.getItem('backgroundColor');
    if (savedColor) {
        document.body.style.backgroundColor = savedColor;
        
        document.querySelectorAll('.color-preset').forEach(btn => {
            if (btn.dataset.color === savedColor) {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });
        
        document.getElementById('customColor').value = savedColor;
    }
}

function setBackgroundColor(color) {
    localStorage.setItem('backgroundColor', color);
    document.body.style.backgroundColor = color;
}

document.addEventListener('DOMContentLoaded', () => {
    applyBackgroundColor();
    
    document.querySelectorAll('.color-preset').forEach(btn => {
        btn.addEventListener('click', () => {
            const color = btn.dataset.color;
            setBackgroundColor(color);
            
            document.querySelectorAll('.color-preset').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            document.getElementById('customColor').value = color;
        });
    });
    
    document.getElementById('applyCustomColor').addEventListener('click', () => {
        const color = document.getElementById('customColor').value;
        setBackgroundColor(color);
        
        document.querySelectorAll('.color-preset').forEach(b => b.classList.remove('active'));
    });
    
    document.getElementById('resetColor').addEventListener('click', () => {
        localStorage.removeItem('backgroundColor');
        document.body.style.backgroundColor = '#ffffff';
        document.getElementById('customColor').value = '#ffffff';
        
        document.querySelectorAll('.color-preset').forEach(btn => {
            if (btn.dataset.color === '#ffffff') {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });
    });
    
    document.querySelector('.dashboard-btn').addEventListener('click', () => {
        window.location.href = '/dashboard';
    });
    
    document.querySelector('.logout-btn').addEventListener("click", async () => {
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
    })})