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

    try {
        const response = await fetch('/api/posts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        });

        if (response.ok) {
            const result = await response.json();
            alert('Post created successfully!');
            window.location.href = '/admin/posts';
        } else if (response.status === 409) {
            alert('Error: Slug already exists. Please choose a different slug.');
        } else {
            alert('Error creating post. Please try again.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Network error. Please try again.');
    }
});

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