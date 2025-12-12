let currentPage = 1;
const limit = 10;

document.addEventListener('DOMContentLoaded', function() {
    loadPages(currentPage);
});

function loadPages(page) {
    fetch(`/api/pages?page=${page}&limit=${limit}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Pages API response:', data); // Debug log
            const tbody = document.getElementById('page-list-body');
            tbody.innerHTML = ''; // Clear existing rows

            // Handle both old array format and new object format
            let pages = data;
            let paginationData = { total: pages.length, page: 1, limit: pages.length, total_pages: 1 };
            if (data.pages) {
                pages = data.pages;
                paginationData = data;
            }

            console.log('Pages array:', pages); // Debug log
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

            updatePagination(paginationData);
        })
        .catch(error => console.error('Error fetching pages:', error));
}

function updatePagination(data) {
    const paginationContainer = document.querySelector('.inline-flex.-space-x-px');
    const showingSpan = document.querySelector('span.text-sm.text-slate-500.dark\\:text-slate-400');

    // Update showing text
    const start = (data.page - 1) * data.limit + 1;
    const end = Math.min(data.page * data.limit, data.total);
    showingSpan.innerHTML = `Showing <span class="font-semibold text-slate-900 dark:text-white">${start}-${end}</span> of <span class="font-semibold text-slate-900 dark:text-white">${data.total}</span>`;

    // Clear existing pagination buttons
    paginationContainer.innerHTML = '';

    // Previous button
    const prevButton = document.createElement('button');
    prevButton.className = 'flex items-center justify-center px-3 h-8 ms-0 leading-tight text-slate-500 bg-white dark:bg-surface-dark border border-e-0 border-border-light dark:border-border-dark rounded-s-lg hover:bg-slate-100 dark:hover:bg-slate-800 dark:text-slate-400';
    prevButton.textContent = 'Previous';
    prevButton.disabled = data.page <= 1;
    if (!prevButton.disabled) {
        prevButton.addEventListener('click', () => {
            currentPage--;
            loadPages(currentPage);
        });
    } else {
        prevButton.classList.add('opacity-50', 'cursor-not-allowed');
    }
    paginationContainer.appendChild(prevButton);

    // Page buttons
    const maxPages = 5;
    let startPage = Math.max(1, data.page - Math.floor(maxPages / 2));
    let endPage = Math.min(data.total_pages, startPage + maxPages - 1);

    if (endPage - startPage + 1 < maxPages) {
        startPage = Math.max(1, endPage - maxPages + 1);
    }

    for (let i = startPage; i <= endPage; i++) {
        const pageButton = document.createElement('button');
        pageButton.className = 'flex items-center justify-center px-3 h-8 leading-tight border border-border-light dark:border-border-dark hover:bg-slate-100 dark:hover:bg-slate-800 dark:text-slate-400';
        pageButton.textContent = i;

        if (i === data.page) {
            pageButton.className += ' text-white bg-primary border-primary hover:bg-blue-600 dark:border-primary dark:text-white';
        } else {
            pageButton.className += ' text-slate-500 bg-white dark:bg-surface-dark';
        }

        pageButton.addEventListener('click', () => {
            currentPage = i;
            loadPages(currentPage);
        });

        paginationContainer.appendChild(pageButton);
    }

    // Next button
    const nextButton = document.createElement('button');
    nextButton.className = 'flex items-center justify-center px-3 h-8 leading-tight text-slate-500 bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-e-lg hover:bg-slate-100 dark:hover:bg-slate-800 dark:text-slate-400';
    nextButton.textContent = 'Next';
    nextButton.disabled = data.page >= data.total_pages;
    if (!nextButton.disabled) {
        nextButton.addEventListener('click', () => {
            currentPage++;
            loadPages(currentPage);
        });
    } else {
        nextButton.classList.add('opacity-50', 'cursor-not-allowed');
    }
    paginationContainer.appendChild(nextButton);
}

function deletePage(id) {
    if (confirm('Are you sure you want to delete this page?')) {
        fetch(`/api/pages/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                // Reload the current page to update pagination
                loadPages(currentPage);
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