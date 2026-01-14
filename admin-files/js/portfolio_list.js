let currentPage = 1;
const limit = 10;

document.addEventListener('DOMContentLoaded', function() {
    loadPortfolioItems(currentPage);
});

function loadPortfolioItems(page) {
    fetch(`/api/portfolio?page=${page}&limit=${limit}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Portfolio API response:', data); // Debug log
            const tbody = document.getElementById('portfolio-list-body');
            tbody.innerHTML = ''; // Clear existing rows

            // Handle both old array format and new object format
            let items = data;
            let paginationData = { total: items.length, page: 1, limit: items.length, total_pages: 1 };
            if (data.items) {
                items = data.items;
                paginationData = data;
            }

            console.log('Items array:', items); // Debug log
            items.forEach(item => {
                const row = document.createElement('tr');
                row.className = 'admin-table-row';
                row.setAttribute('data-portfolio-id', item.id);

                const thumbnailHtml = item.showcase_image ?
                    `<img src="${item.showcase_image}" alt="${item.title}" class="table-thumbnail">` :
                    '<span class="text-muted">No image</span>';

                row.innerHTML = `
                <td class="table-cell-title">${item.title}</td>
                <td>${thumbnailHtml}</td>
                <td>${item.sort_order}</td>
                <td class="text-right">
                    <div class="table-cell-actions">
                        <a href="/admin/portfolio/${item.id}/edit">
                            <button class="table-action-btn" title="Edit">
                                <span class="material-symbols-outlined">edit</span>
                            </button>
                        </a>
                        <button onclick="deletePortfolioItem(${item.id})" class="table-action-btn table-action-btn-danger" title="Delete">
                            <span class="material-symbols-outlined">delete</span>
                        </button>
                    </div>
                </td>
                `;
                tbody.appendChild(row);
            });

            updatePagination(paginationData);
        })
        .catch(error => console.error('Error fetching portfolio items:', error));
}

function updatePagination(data) {
    const paginationContainer = document.querySelector('.pagination-buttons');
    const showingSpan = document.querySelector('.pagination-info');

    // Update showing text
    const start = (data.page - 1) * data.limit + 1;
    const end = Math.min(data.page * data.limit, data.total);
    showingSpan.innerHTML = `Showing <span class="pagination-highlight">${start}-${end}</span> of <span class="pagination-highlight">${data.total}</span>`;

    // Clear existing pagination buttons
    paginationContainer.innerHTML = '';

    // Previous button
    const prevButton = document.createElement('button');
    prevButton.className = 'pagination-btn';
    prevButton.textContent = 'Previous';
    prevButton.disabled = data.page <= 1;
    if (!prevButton.disabled) {
        prevButton.addEventListener('click', () => {
            currentPage--;
            loadPortfolioItems(currentPage);
        });
    } else {
        prevButton.classList.add('pagination-btn-disabled');
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
        pageButton.className = i === data.page ? 'pagination-btn pagination-btn-active' : 'pagination-btn';
        pageButton.textContent = i;

        pageButton.addEventListener('click', () => {
            currentPage = i;
            loadPortfolioItems(currentPage);
        });

        paginationContainer.appendChild(pageButton);
    }

    // Next button
    const nextButton = document.createElement('button');
    nextButton.className = 'pagination-btn';
    nextButton.textContent = 'Next';
    nextButton.disabled = data.page >= data.total_pages;
    if (!nextButton.disabled) {
        nextButton.addEventListener('click', () => {
            currentPage++;
            loadPortfolioItems(currentPage);
        });
    } else {
        nextButton.classList.add('pagination-btn-disabled');
    }
    paginationContainer.appendChild(nextButton);
}

function deletePortfolioItem(id) {
    if (confirm('Are you sure you want to delete this portfolio item?')) {
        fetch(`/api/portfolio/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                // Reload the current page to update pagination
                loadPortfolioItems(currentPage);
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