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
const formattingToolbar = document.getElementById('formattingToolbar');

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
        formattingToolbar.style.display = 'flex';
    } else {
        editorMode.classList.remove('active');
        viewerMode.classList.add('active');
        toggleButton.textContent = 'Editar';
        formattingToolbar.style.display = 'none';
        
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

function insertMarkdown(format) {
    const start = editor.selectionStart;
    const end = editor.selectionEnd;
    const selectedText = editor.value.substring(start, end);
    const beforeText = editor.value.substring(0, start);
    const afterText = editor.value.substring(end);
    
    let newText = '';
    let cursorOffset = 0;
    
    switch (format) {
        case 'h1':
            if (start === 0 || beforeText.endsWith('\n')) {
                newText = `# ${selectedText || 'Título'}`;
                cursorOffset = selectedText ? 0 : -6; // Move cursor before "Título"
            } else {
                newText = `\n# ${selectedText || 'Título'}`;
                cursorOffset = selectedText ? 0 : -6;
            }
            break;
            
        case 'h2':
            if (start === 0 || beforeText.endsWith('\n')) {
                newText = `## ${selectedText || 'Subtítulo'}`;
                cursorOffset = selectedText ? 0 : -10;
            } else {
                newText = `\n## ${selectedText || 'Subtítulo'}`;
                cursorOffset = selectedText ? 0 : -10;
            }
            break;
            
        case 'h3':
            if (start === 0 || beforeText.endsWith('\n')) {
                newText = `### ${selectedText || 'Título 3'}`;
                cursorOffset = selectedText ? 0 : -9;
            } else {
                newText = `\n### ${selectedText || 'Título 3'}`;
                cursorOffset = selectedText ? 0 : -9;
            }
            break;
            
        case 'bold':
            newText = `**${selectedText || 'texto en negrita'}**`;
            cursorOffset = selectedText ? 0 : -19;
            break;
            
        case 'italic':
            newText = `*${selectedText || 'texto en cursiva'}*`;
            cursorOffset = selectedText ? 0 : -18;
            break;
            
        case 'strikethrough':
            newText = `~~${selectedText || 'texto tachado'}~~`;
            cursorOffset = selectedText ? 0 : -16;
            break;
            
        case 'code':
            newText = `\`${selectedText || 'código'}\``;
            cursorOffset = selectedText ? 0 : -7;
            break;
            
        case 'codeblock':
            if (start === 0 || beforeText.endsWith('\n')) {
                newText = `\`\`\`\n${selectedText || 'código'}\n\`\`\``;
                cursorOffset = selectedText ? 0 : -11;
            } else {
                newText = `\n\`\`\`\n${selectedText || 'código'}\n\`\`\``;
                cursorOffset = selectedText ? 0 : -11;
            }
            break;
            
        case 'ul':
            if (selectedText) {
                const lines = selectedText.split('\n');
                newText = lines.map(line => `- ${line}`).join('\n');
            } else {
                if (start === 0 || beforeText.endsWith('\n')) {
                    newText = '- elemento de lista';
                } else {
                    newText = '\n- elemento de lista';
                }
                cursorOffset = -18;
            }
            break;
            
        case 'ol':
            if (selectedText) {
                const lines = selectedText.split('\n');
                newText = lines.map((line, i) => `${i + 1}. ${line}`).join('\n');
            } else {
                if (start === 0 || beforeText.endsWith('\n')) {
                    newText = '1. elemento de lista';
                } else {
                    newText = '\n1. elemento de lista';
                }
                cursorOffset = -19;
            }
            break;
            
        case 'quote':
            if (selectedText) {
                const lines = selectedText.split('\n');
                newText = lines.map(line => `> ${line}`).join('\n');
            } else {
                if (start === 0 || beforeText.endsWith('\n')) {
                    newText = '> cita';
                } else {
                    newText = '\n> cita';
                }
                cursorOffset = -4;
            }
            break;
            
        case 'link':
            newText = `[${selectedText || 'texto del enlace'}](url)`;
            cursorOffset = selectedText ? -5 : -22;
            break;
            
        case 'image':
            newText = `![${selectedText || 'descripción'}](url)`;
            cursorOffset = selectedText ? -5 : -19;
            break;
            
        case 'hr':
            if (start === 0 || beforeText.endsWith('\n')) {
                newText = '---';
            } else {
                newText = '\n---';
            }
            cursorOffset = 0;
            break;
    }
    
    editor.value = beforeText + newText + afterText;
    
    const newCursorPos = start + newText.length + cursorOffset;
    editor.setSelectionRange(newCursorPos, newCursorPos);
    editor.focus();
    
    handleInput();
}

document.querySelectorAll('.format-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        e.preventDefault();
        const format = btn.dataset.format;
        insertMarkdown(format);
    });
});

init();


