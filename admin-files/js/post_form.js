document.getElementById('post-form').addEventListener('submit', async function(e) {
    e.preventDefault();

    const title = document.getElementById('title').value.trim();
    const slug = document.getElementById('slug').value.trim();
    const tags = document.getElementById('tags').value.trim();
    const content = document.getElementById('content').value.trim();
    const published = document.getElementById('published').checked;

    // Basic validation
    if (!title || !slug || !content) {
        alert('Please fill in all required fields: Title, Slug, and Content.');
        return;
    }

    const postData = {
        title: title,
        slug: slug,
        tags: tags,
        content: content,
        published: published
    };

    const isEdit = window.location.pathname.includes('/edit');
    const url = isEdit ? `/api/posts/${postId}` : '/api/posts';
    const method = isEdit ? 'PUT' : 'POST';
    const successMessage = isEdit ? 'Post updated successfully!' : 'Post created successfully!';

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        });

        if (response.ok) {
            alert(successMessage);
            window.location.href = '/admin/posts';
        } else if (response.status === 409) {
            alert('Error: Slug already exists. Please choose a different slug.');
        } else {
            alert(`Error ${isEdit ? 'updating' : 'creating'} post. Please try again.`);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Network error. Please try again.');
    }
});

// Check if editing
let postId = null;
const pathMatch = window.location.pathname.match(/^\/admin\/posts\/(\d+)\/edit$/);
if (pathMatch) {
    postId = pathMatch[1];
    // Fetch post data
    fetch(`/api/posts/${postId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Post not found');
            }
            return response.json();
        })
        .then(post => {
            document.getElementById('title').value = post.title || '';
            document.getElementById('slug').value = post.slug || '';
            document.getElementById('tags').value = post.tags || '';
            document.getElementById('content').value = post.content || '';
            document.getElementById('published').checked = post.published || false;
            document.getElementById('slug').dataset.original = post.slug || '';
        })
        .catch(error => {
            console.error('Error loading post:', error);
            alert('Error loading post data.');
        });
}

// Auto-slugify title
document.getElementById('title').addEventListener('input', function() {
    const title = this.value.trim();
    const slugField = document.getElementById('slug');
    if (slugField.value === '' || slugField.value === slugify(slugField.dataset.original || '')) {
        slugField.value = slugify(title);
        slugField.dataset.original = title;
    }
});

function slugify(text) {
    return text
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .replace(/[\s_-]+/g, '-')
        .replace(/^-+|-+$/g, '');
}