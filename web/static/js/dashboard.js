document.addEventListener('DOMContentLoaded', async function() {
    const response = await fetch('/api/me', { credentials: 'include' });
    if (!response.ok) {
        window.location.href = '/login';
        return;
    }

    await loadNotes();
});

const $logoutBtn = document.querySelector(".logout-btn")
const $settingsBtn = document.querySelector(".settings-btn")
const $profileBtn = document.querySelector(".profile-btn")
const $dashboardBtn = document.querySelector(".dashboard-btn")
const $createNoteBtn = document.querySelector(".create-note-btn")
const $uploadNoteBtn = document.querySelector(".upload-note-btn")
const $mdFileInput = document.getElementById("md-file-input")
const $notesContainer = document.getElementById("notes-container")
const $createNoteOverlay = document.querySelector(".create-note-overlay")
const $noteTitleInput = document.querySelector(".note-title-input")
const $createNoteButton = document.querySelector(".create-note-button")
const $cancelNoteButton = document.querySelector(".cancel-note-button")

async function loadNotes() {
    try {
        const response = await fetch('/api/me/notes?limit=100&offset=0', {
            credentials: 'include'
        });
        
        if (!response.ok) {
            throw new Error('Failed to fetch notes');
        }
        
        const data = await response.json();
        displayNotes(data.data || []);
    } catch (error) {
        console.error('Error loading notes:', error);
        $notesContainer.innerHTML = '<div class="empty-state"><h2>Error</h2><p>No se pudieron cargar las notas</p></div>';
    }
}

function displayNotes(notes) {
    if (notes.length === 0) {
        $notesContainer.innerHTML = `
            <div class="empty-state">
                <h2>No hay notas</h2>
                <p>Crea tu primera nota para comenzar</p>
            </div>
        `;
        return;
    }
    
    $notesContainer.innerHTML = notes.map(note => `
        <div class="note-card" data-note-id="${note.id}">
            <div class="note-card-title">${escapeHtml(note.title)}</div>
            <div class="note-card-preview">${escapeHtml(note.content.substring(0, 150))}${note.content.length > 150 ? '...' : ''}</div>
            <div class="note-card-date">${formatDate(note.updatedAt)}</div>
            <div class="note-card-actions">
                <button class="edit-note-btn" data-note-id="${note.id}" title="Editar nota">
                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" pointer-events="none">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>
                <button class="delete-note-btn" data-note-id="${note.id}" title="Eliminar nota">
                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" pointer-events="none">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                        <path d="M10 11v6M14 11v6"/>
                        <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
                    </svg>
                </button>
            </div>
        </div>
    `).join('');
    
    document.querySelectorAll('.note-card').forEach(card => {
        card.addEventListener('click', (e) => {
            const noteId = card.dataset.noteId;
            window.location.href = `/notes/${noteId}`;
        });
    });
    
    document.querySelectorAll('.edit-note-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.stopPropagation();
            const noteId = e.target.dataset.noteId;
            window.location.href = `/notes/${noteId}`;
        });
    });
    
    document.querySelectorAll('.delete-note-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            e.stopPropagation();
            const noteId = e.target.dataset.noteId;
            await deleteNote(noteId);
        });
    });
}

async function deleteNote(noteId) {
    if (!confirm('¿Estás seguro de que quieres eliminar esta nota?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/notes/${noteId}`, {
            method: 'DELETE',
            credentials: 'include'
        });
        
        if (response.ok) {
            const noteCard = document.querySelector(`[data-note-id="${noteId}"]`);
            if (noteCard) {
                noteCard.remove();
            }
            
            if (document.querySelectorAll('.note-card').length === 0) {
                displayNotes([]);
            }
        } else {
            alert('Error al eliminar la nota');
        }
    } catch (error) {
        console.error('Error deleting note:', error);
        alert('Error al eliminar la nota');
    }
}

function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
        return 'Hoy';
    } else if (diffDays === 1) {
        return 'Ayer';
    } else if (diffDays < 7) {
        return `Hace ${diffDays} días`;
    } else {
        return date.toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' });
    }
}

$logoutBtn.addEventListener("click", async () => {
    try {
        const response = await fetch("/api/auth/logout", {
            method: "DELETE",
            credentials: 'include'
        })

        if (response.ok) {
            window.location.href = "/"
        } else {
            const errorText = await response.text()
            alert(errorText)
        }
    } catch (error) {
        alert("Error al conectar con el servidor. Verifica tu conexión a internet")
        console.error(error)
    }
})

$settingsBtn.addEventListener("click", () => {
    window.location.href = "/settings"
})

$profileBtn.addEventListener("click", () => {
    window.location.href = "/me"
})

$dashboardBtn.addEventListener("click", () => {
    window.location.href = "/dashboard"
})

$createNoteBtn.addEventListener("click", () => {
    $createNoteOverlay.style.display = "flex"
    $noteTitleInput.value = ""
    $noteTitleInput.focus()
})

$cancelNoteButton.addEventListener("click", () => {
    $createNoteOverlay.style.display = "none"
})

$createNoteOverlay.addEventListener("click", (e) => {
    if (e.target === $createNoteOverlay) {
        $createNoteOverlay.style.display = "none"
    }
})

$createNoteButton.addEventListener("click", async () => {
    const title = $noteTitleInput.value.trim()
    
    if (!title) {
        alert("Por favor ingresa un título")
        return
    }
    
    try {
        const response = await fetch("/api/notes", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({
                title: title,
                content: ""
            })
        })
        
        if (response.ok) {
            const newNote = await response.json()
            $createNoteOverlay.style.display = "none"
            // Redirect to edit the new note
            window.location.href = `/notes/${newNote.id}`
        } else {
            const errorText = await response.text()
            alert("Error al crear la nota: " + errorText)
        }
    } catch (error) {
        console.error("Error creating note:", error)
        alert("Error al crear la nota")
    }
})

// Allow pressing Enter to create note
$noteTitleInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") {
        $createNoteButton.click()
    }
})

$uploadNoteBtn.addEventListener("click", () => {
    $mdFileInput.click()
})

$mdFileInput.addEventListener("change", async (e) => {
    const file = e.target.files[0]
    if (!file) return
    
    if (!file.name.endsWith('.md')) {
        alert("Por favor selecciona un archivo .md")
        return
    }
    
    try {
        const content = await file.text()
        
        const response = await fetch("/api/notes", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({
                title: "nota subida",
                content: content
            })
        })
        
        if (response.ok) {
            const newNote = await response.json()
            alert("Nota subida exitosamente")
            // Reload notes to show the new one
            await loadNotes()
        } else {
            const errorText = await response.text()
            alert("Error al subir la nota: " + errorText)
        }
    } catch (error) {
        console.error("Error uploading file:", error)
        alert("Error al leer o subir el archivo")
    }
    
    e.target.value = ''
})
