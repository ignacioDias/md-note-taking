let debounceTimer;
let noteId;
let isEditorMode = true;

const titleInput = document.getElementById('title');
const editor = document.getElementById('editor');
const saveStatus = document.getElementById('save-status');
const toggleButton = document.getElementById('toggleMode');
const editorMode = document.getElementById('editorMode');
const viewerMode = document.getElementById('viewerMode');
const viewerTitle = document.getElementById('viewerTitle');
const viewerContent = document.getElementById('viewerContent');
const comeBackButton = document.querySelector('.come-back')

comeBackButton.addEventListener("click", () => {
    window.location.href = '/dashboard';
})

if (typeof marked !== 'undefined') {
    marked.setOptions({
        breaks: true,
        gfm: true
    });
}

async function init() {
    const note = await loadNote();
    if (note) {
        noteId = note.id;
        titleInput.value = note.title || '';
        editor.value = note.content || '';
    }
}

async function loadNote() {
    try {
        const id = window.location.pathname.split("/")[2];
        const res = await fetch(`/api/notes/${id}`, {
            credentials: 'include'
        });
        
        if (!res.ok) {
            throw new Error('Failed to load note');
        }
        
        const note = await res.json();
        return note;
    } catch (error) {
        console.error('Error loading note:', error);
        saveStatus.textContent = '❌ Error loading note';
        return null;
    }
}

async function saveNote() {
    if (!noteId) return;
    
    saveStatus.textContent = 'Guardando...';
    
    try {
        const response = await fetch(`/api/notes/${noteId}`, {
            method: 'PUT',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                title: titleInput.value,
                content: editor.value
            })
        });
        
        if (response.ok) {
            saveStatus.textContent = '✓ Guardado';
        } else {
            saveStatus.textContent = '❌ Error guardando';
        }
    } catch (error) {
        console.error('Error guardando note:', error);
        saveStatus.textContent = '❌ Error guardando';
    }
}

function handleInput() {
    clearTimeout(debounceTimer);
    saveStatus.textContent = 'Editando...';
    debounceTimer = setTimeout(saveNote, 1000);
}

function toggleMode() {
    isEditorMode = !isEditorMode;
    
    if (isEditorMode) {
        editorMode.classList.add('active');
        viewerMode.classList.remove('active');
        toggleButton.textContent = 'Preview';
    } else {
        editorMode.classList.remove('active');
        viewerMode.classList.add('active');
        toggleButton.textContent = 'Editar';
        
        viewerTitle.textContent = titleInput.value || 'Untitled';
        if (typeof marked !== 'undefined') {
            viewerContent.innerHTML = marked.parse(editor.value || '*No content*');
        } else {
            viewerContent.textContent = editor.value || 'No content';
        }
    }
}

titleInput.addEventListener('input', handleInput);
editor.addEventListener('input', handleInput);
toggleButton.addEventListener('click', toggleMode);

window.addEventListener('beforeunload', (e) => {
    if (saveStatus.textContent === 'Editando...') {
        clearTimeout(debounceTimer);
        saveNote();
    }
});

init();


