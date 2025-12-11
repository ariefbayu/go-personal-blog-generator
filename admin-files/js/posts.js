document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/posts')
        .then(response => response.json())
        .then(posts => {
            const tbody = document.getElementById('post-list-body');
            posts.forEach(post => {
                const row = document.createElement('tr');
                row.className = 'bg-surface-light dark:bg-surface-dark hover:bg-slate-50 dark:hover:bg-[#1e2a36] transition-colors';
                row.innerHTML = `

                <td
                    class="px-6 py-4 font-medium text-slate-900 dark:text-white whitespace-nowrap">
                    ${post.title}
                </td>
                <td class="px-6 py-4">
                    <span
                        class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-green-500"></span>
                        ${post.published ? `
                            <span
                                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">
                                <span class="w-1.5 h-1.5 rounded-full bg-green-500"></span>
                                Published
                            </span>` : `
                            <span
                                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-slate-100 text-slate-800 dark:bg-slate-700/50 dark:text-slate-300">
                                <span class="w-1.5 h-1.5 rounded-full bg-slate-400"></span>
                                Draft
                            </span>`}
                    </span>
                </td>
                <td class="px-6 py-4 text-right">
                    <div class="flex items-center justify-end gap-2">
                        <a href="/admin/posts/${post.id}/edit">
                            <button
                                class="p-2 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                                title="Edit">
                                <span class="material-symbols-outlined text-[20px]">edit</span>
                            </button>
                        </a>
                        <button onclick="deletePost(${post.id})"
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
        .catch(error => console.error('Error fetching posts:', error));
});

function deletePost(id) {
    if (confirm('Are you sure you want to delete this post?')) {
        // Placeholder for delete functionality
        alert('Delete functionality not implemented yet');
    }
}
