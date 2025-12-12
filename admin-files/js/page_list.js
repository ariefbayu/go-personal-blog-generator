document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/pages')
        .then(response => response.json())
        .then(pages => {
            const tbody = document.getElementById('page-list-body');
            pages.forEach(page => {
                const row = document.createElement('tr');
                row.className = 'bg-surface-light dark:bg-surface-dark hover:bg-slate-50 dark:hover:bg-[#1e2a36] transition-colors';
                row.setAttribute('data-page-id', page.id);

                const navStatus = page.show_in_nav ?
                    '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300">Yes</span>' :
                    '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300">No</span>';

                row.innerHTML = `
                <td class="px-6 py-4 font-medium text-slate-900 dark:text-white whitespace-nowrap">
                    ${page.title}
                </td>
                <td class="px-6 py-4 text-slate-500 dark:text-slate-400">
                    ${page.slug}
                </td>
                <td class="px-6 py-4">
                    ${navStatus}
                </td>
                <td class="px-6 py-4">
                    ${page.sort_order}
                </td>
                <td class="px-6 py-4 text-right">
                    <div class="flex items-center justify-end gap-2">
                        <a href="/admin/pages/${page.id}/edit">
                            <button
                                class="p-2 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                                title="Edit">
                                <span class="material-symbols-outlined text-[20px]">edit</span>
                            </button>
                        </a>
                        <button onclick="deletePage(${page.id})"
                            class="p-2 text-slate-500 dark:text-slate-400 hover:text-red-500 dark:hover:text-red-500 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                            title="Delete">
                            <span class="material-symbols-outlined text-[20px]">delete</span>
                        </button>
                    </div>
                </td>
                `;
                tbody.appendChild(row);
            });
        })
        .catch(error => console.error('Error fetching pages:', error));
});

function deletePage(id) {
    if (confirm('Are you sure you want to delete this page?')) {
        fetch(`/api/pages/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                // Remove the row from the table
                const row = document.querySelector(`[data-page-id="${id}"]`);
                if (row) {
                    row.remove();
                }
                alert('Page deleted successfully');
            } else if (response.status === 404) {
                alert('Page not found');
            } else {
                alert('Error deleting page');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Network error. Please try again.');
        });
    }
}