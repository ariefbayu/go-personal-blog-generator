document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/portfolio')
        .then(response => response.json())
        .then(portfolioItems => {
            const tbody = document.getElementById('portfolio-list-body');
            portfolioItems.forEach(item => {
                const row = document.createElement('tr');
                row.className = 'bg-surface-light dark:bg-surface-dark hover:bg-slate-50 dark:hover:bg-[#1e2a36] transition-colors';
                row.setAttribute('data-portfolio-id', item.id);

                const thumbnailHtml = item.showcase_image ?
                    `<img src="${item.showcase_image}" alt="${item.title}" class="w-16 h-16 object-cover rounded border">` :
                    '<span class="text-slate-400">No image</span>';

                row.innerHTML = `
                <td class="px-6 py-4 font-medium text-slate-900 dark:text-white whitespace-nowrap">
                    ${item.title}
                </td>
                <td class="px-6 py-4">
                    ${thumbnailHtml}
                </td>
                <td class="px-6 py-4">
                    ${item.sort_order}
                </td>
                <td class="px-6 py-4 text-right">
                    <div class="flex items-center justify-end gap-2">
                        <a href="/admin/portfolio/${item.id}/edit">
                            <button
                                class="p-2 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                                title="Edit">
                                <span class="material-symbols-outlined text-[20px]">edit</span>
                            </button>
                        </a>
                        <button onclick="deletePortfolioItem(${item.id})"
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
        .catch(error => console.error('Error fetching portfolio items:', error));
});

function deletePortfolioItem(id) {
    if (confirm('Are you sure you want to delete this portfolio item?')) {
        fetch(`/api/portfolio/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                // Remove the row from the table
                const row = document.querySelector(`[data-portfolio-id="${id}"]`);
                if (row) {
                    row.remove();
                }
                alert('Portfolio item deleted successfully');
            } else if (response.status === 404) {
                alert('Portfolio item not found');
            } else {
                alert('Error deleting portfolio item');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Network error. Please try again.');
        });
    }
}