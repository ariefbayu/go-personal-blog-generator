let currentFilePath = null;

document.addEventListener('DOMContentLoaded', function() {
    loadFileTree();
});

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
            document.getElementById('file-tree').innerHTML = '<p>Error loading file tree.</p>';
        });
}

function buildTree(nodes, container, path) {
    const ul = document.createElement('ul');
    ul.className = 'file-tree-list';
    nodes.forEach(node => {
        const li = document.createElement('li');
        li.className = 'file-tree-item';
        const fullPath = path ? path + '/' + node.name : node.name;
        li.innerHTML = `
            <span class="file-tree-toggle">${node.type === 'dir' ? '▶' : ''}</span>
            <span class="file-tree-icon material-symbols-outlined">${node.type === 'dir' ? 'folder' : 'description'}</span>
            <span class="file-tree-name ${node.editable ? 'editable' : ''}" data-path="${fullPath}" data-type="${node.type}">${node.name}</span>
        `;
        if (node.type === 'dir') {
            li.classList.add('file-tree-dir');
            const childUl = document.createElement('ul');
            childUl.className = 'file-tree-children';
            childUl.style.display = 'none';
            buildTree(node.children, childUl, fullPath);
            li.appendChild(childUl);
            li.querySelector('.file-tree-toggle').addEventListener('click', function() {
                toggleDir(this, childUl);
            });
        } else {
            li.querySelector('.file-tree-name').addEventListener('click', function() {
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
    const textarea = document.getElementById('file-content');
    textarea.readOnly = false;
    document.getElementById('save-btn').disabled = false;

    fetch(`/api/settings/templates/content?path=${encodeURIComponent(path)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.text();
        })
        .then(content => {
            textarea.value = content;
        })
        .catch(error => {
            console.error('Error loading file content:', error);
            textarea.value = 'Error loading file content.';
        });
}

document.getElementById('save-btn').addEventListener('click', function() {
    if (!currentFilePath) return;

    const content = document.getElementById('file-content').value;
    const status = document.getElementById('save-status');

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
        status.style.color = 'green';
        setTimeout(() => { status.textContent = ''; }, 3000);
    })
    .catch(error => {
        console.error('Error saving file:', error);
        status.textContent = 'Error saving file.';
        status.style.color = 'red';
        setTimeout(() => { status.textContent = ''; }, 3000);
    });
});