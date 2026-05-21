document.addEventListener('DOMContentLoaded', function() {
    loadSettings();

    document.getElementById('settingsForm').addEventListener('submit', function(e) {
        e.preventDefault();
        saveSettings();
    });
});

async function loadSettings() {
    try {
        const response = await fetch(`${window.ROOT_PREFIX || ''}/api/settings`);
        const settings = await response.json();
        document.getElementById('siteName').value = settings.site_name;
        document.getElementById('showPortfolioMenu').checked = settings.show_portfolio_menu;
        document.getElementById('showPostsMenu').checked = settings.show_posts_menu;

        let menuOrder = JSON.parse(settings.menu_order || '[]');
        if (menuOrder.length === 0) {
            menuOrder = ["posts", "portfolio", "pages"];
        }
        const list = document.getElementById('menuOrderList');
        list.innerHTML = '';
        menuOrder.forEach(item => {
            const li = document.createElement('li');
            li.className = 'drag-item';
            li.draggable = true;
            li.dataset.item = item;
            li.innerHTML = `<span class="material-symbols-outlined">drag_indicator</span>${item}`;
            li.addEventListener('dragstart', handleDragStart);
            li.addEventListener('dragend', handleDragEnd);
            li.addEventListener('dragover', handleDragOver);
            li.addEventListener('drop', handleDrop);
            list.appendChild(li);
        });
    } catch (error) {
        console.error('Error loading settings:', error);
    }
}

async function saveSettings() {
    const list = document.getElementById('menuOrderList');
    const items = Array.from(list.children).map(li => li.dataset.item);
    const menuOrder = JSON.stringify(items);

    const settings = {
        site_name: document.getElementById('siteName').value,
        show_portfolio_menu: document.getElementById('showPortfolioMenu').checked,
        show_posts_menu: document.getElementById('showPostsMenu').checked,
        menu_order: menuOrder
    };

    try {
        const response = await fetch(`${window.ROOT_PREFIX || ''}/api/settings`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(settings)
        });
        const result = await response.json();
        alert(result.message);
    } catch (error) {
        console.error('Error saving settings:', error);
        alert('Error saving settings');
    }
}

let draggedElement = null;

function handleDragStart(e) {
    draggedElement = e.target;
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/html', e.target.outerHTML);
    e.target.style.opacity = '0.5';
}

function handleDragOver(e) {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    return false;
}

function handleDragEnd(e) {
    e.target.style.opacity = '1';
}

function handleDrop(e) {
    e.preventDefault();
    if (draggedElement !== e.target && e.target.tagName === 'LI') {
        const list = document.getElementById('menuOrderList');
        const allItems = Array.from(list.children);
        const draggedIndex = allItems.indexOf(draggedElement);
        const targetIndex = allItems.indexOf(e.target);

        if (draggedIndex < targetIndex) {
            list.insertBefore(draggedElement, e.target.nextSibling);
        } else {
            list.insertBefore(draggedElement, e.target);
        }
    }
    draggedElement = null;
    return false;
}
