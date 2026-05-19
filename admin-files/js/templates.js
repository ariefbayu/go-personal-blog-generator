let currentFilePath = null;
let editor = null;

document.addEventListener('DOMContentLoaded', function() {
    initEditor();
    loadFileTree();
});

function initEditor() {
    const textarea = document.getElementById('file-content');
    editor = CodeMirror.fromTextArea(textarea, {
        lineNumbers: true,
        mode: "htmlmixed",
        theme: "material-ocean",
        tabSize: 4,
        indentUnit: 4,
        lineWrapping: true,
        viewportMargin: Infinity
    });
}

function loadFileTree() {
    fetch('/api/settings/templates')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const treeContainer = document.getElementById('file-tree');
            treeContainer.innerHTML = '';
            buildTree(data, treeContainer, '');
        })
        .catch(error => {
            console.error('Error loading file tree:', error);
            document.getElementById('file-tree').innerHTML = '<p class="text-danger">Error loading file tree.</p>';
        });
}

function buildTree(nodes, container, path) {
    const ul = document.createElement('ul');
    ul.className = 'file-tree-list';
    nodes.forEach(node => {
        const li = document.createElement('li');
        li.className = 'file-tree-item-wrapper';
        const fullPath = path ? path + '/' + node.name : node.name;
        
        const item = document.createElement('div');
        item.className = 'file-tree-item';
        item.innerHTML = `
            <span class="file-tree-toggle">${node.type === 'dir' ? '▶' : ''}</span>
            <span class="file-tree-icon material-symbols-outlined">${node.type === 'dir' ? 'folder' : 'description'}</span>
            <span class="file-tree-name ${node.editable ? 'editable' : ''}" data-path="${fullPath}" data-type="${node.type}">${node.name}</span>
        `;
        li.appendChild(item);

        if (node.type === 'dir') {
            const childUl = document.createElement('ul');
            childUl.className = 'file-tree-children';
            childUl.style.display = 'none';
            buildTree(node.children, childUl, fullPath);
            li.appendChild(childUl);
            
            item.querySelector('.file-tree-toggle').addEventListener('click', (e) => {
                e.stopPropagation();
                toggleDir(item.querySelector('.file-tree-toggle'), childUl);
            });
            item.addEventListener('click', () => {
                toggleDir(item.querySelector('.file-tree-toggle'), childUl);
            });
        } else {
            item.addEventListener('click', () => {
                selectFile(fullPath);
            });
        }
        ul.appendChild(li);
    });
    container.appendChild(ul);
}

function toggleDir(toggle, childUl) {
    if (childUl.style.display === 'none') {
        childUl.style.display = 'block';
        toggle.textContent = '▼';
    } else {
        childUl.style.display = 'none';
        toggle.textContent = '▶';
    }
}

function selectFile(path) {
    currentFilePath = path;
    document.getElementById('editor-title').textContent = `Editing: ${path}`;
    document.getElementById('save-btn').disabled = false;

    // Determine mode based on extension
    const ext = path.split('.').pop().toLowerCase();
    let mode = "htmlmixed";
    if (ext === "css") mode = "css";
    if (ext === "js") mode = "javascript";
    if (ext === "json") mode = "application/json";
    if (ext === "md") mode = "markdown";
    
    editor.setOption("mode", mode);

    fetch(`/api/settings/templates/content?path=${encodeURIComponent(path)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.text();
        })
        .then(content => {
            editor.setValue(content);
        })
        .catch(error => {
            console.error('Error loading file content:', error);
            editor.setValue('Error loading file content.');
        });
}

document.getElementById('save-btn').addEventListener('click', function() {
    if (!currentFilePath) return;

    const content = editor.getValue();
    const status = document.getElementById('save-status');
    const btn = this;
    const originalText = btn.innerHTML;

    btn.disabled = true;
    btn.innerHTML = '<span class="material-symbols-outlined">sync</span><span>Saving...</span>';

    fetch('/api/settings/templates/save', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ path: currentFilePath, content: content }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        status.textContent = 'File saved successfully!';
        status.className = 'save-status-msg text-success';
        setTimeout(() => { status.textContent = ''; }, 3000);
    })
    .catch(error => {
        console.error('Error saving file:', error);
        status.textContent = 'Error saving file.';
        status.className = 'save-status-msg text-danger';
        setTimeout(() => { status.textContent = ''; }, 3000);
    })
    .finally(() => {
        btn.disabled = false;
        btn.innerHTML = originalText;
    });
});
