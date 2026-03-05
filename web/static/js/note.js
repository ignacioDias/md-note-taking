let debounceTimer;

editor.addEventListener('input', function() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(async () => {
        await fetch(`/api/notes/${noteId}`, {
            method: 'PUT',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                title: titleInput.value,
                content: editor.value
            })
        });
    }, 1000); // guarda 1 segundo después de que para de escribir
});