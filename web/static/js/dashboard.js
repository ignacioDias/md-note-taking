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
                <button class="edit-note-btn" data-note-id="${note.id}" title="Editar nota">✏️</button>
                <button class="delete-note-btn" data-note-id="${note.id}" title="Eliminar nota">🗑️</button>
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
