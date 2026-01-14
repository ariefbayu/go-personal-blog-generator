let currentPage = 1;
const limit = 10;

document.addEventListener('DOMContentLoaded', function() {
    loadPosts(currentPage);
});

function loadPosts(page) {
    fetch(`/api/posts?page=${page}&limit=${limit}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('API response:', data); // Debug log
            const tbody = document.getElementById('post-list-body');
            tbody.innerHTML = ''; // Clear existing rows

            // Handle both old array format and new object format
            let posts = data;
            let paginationData = { total: posts.length, page: 1, limit: posts.length, total_pages: 1 };
            if (data.posts) {
                posts = data.posts;
                paginationData = data;
            }

            console.log('Posts array:', posts); // Debug log
            posts.forEach(post => {
                const row = document.createElement('tr');
                row.className = 'admin-table-row';
                row.setAttribute('data-post-id', post.id);
                row.innerHTML = `
                <td class="table-cell-title">${post.title}</td>
                <td>
                    ${post.published ? `
                        <span class="badge-status badge-success">
                            <span class="badge-dot"></span>
                            Published
                        </span>` : `
                        <span class="badge-status badge-secondary">
                            <span class="badge-dot"></span>
                            Draft
                        </span>`}
                </td>
                <td class="text-right">
                    <div class="table-cell-actions">
                        <a href="/admin/posts/${post.id}/edit">
                            <button class="table-action-btn" title="Edit">
                                <span class="material-symbols-outlined">edit</span>
                            </button>
                        </a>
                        <button onclick="deletePost(${post.id})" class="table-action-btn table-action-btn-danger" title="Delete">
                            <span class="material-symbols-outlined">delete</span>
                        </button>
                    </div>
                </td>
                `;
                tbody.appendChild(row);
            });

            updatePagination(paginationData);
        })
        .catch(error => console.error('Error fetching posts:', error));
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
            loadPosts(currentPage);
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
            loadPosts(currentPage);
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
            loadPosts(currentPage);
        });
    } else {
        nextButton.classList.add('pagination-btn-disabled');
    }
    paginationContainer.appendChild(nextButton);
}

function deletePost(id) {
    if (confirm('Are you sure you want to delete this post?')) {
        fetch(`/api/posts/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                // Reload the current page to update pagination
                loadPosts(currentPage);
                alert('Post deleted successfully');
            } else if (response.status === 404) {
                alert('Post not found');
            } else {
                alert('Error deleting post');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Network error. Please try again.');
        });
    }
}
