document.getElementById('post-form').addEventListener('submit', async function(e) {

// Check if editing
let postId = null;
const pathMatch = window.location.pathname.match(/^\/admin\/posts\/(\d+)\/edit$/);
if (pathMatch) {
    postId = pathMatch[1];
}

    e.preventDefault();

    const title = document.getElementById('title').value.trim();
    const slug = document.getElementById('slug').value.trim();
    const tags = document.getElementById('tags').value.trim();

    // Get content from EasyMDE if editor is initialized, otherwise from textarea
    let content;
    if (easyMDE) {
        content = easyMDE.value().trim();
    } else {
        content = document.getElementById('content').value.trim();
    }
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

    console.log('Sending post data:', postData);
    console.log('JSON string:', JSON.stringify(postData));

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

// Initialize EasyMDE editor
let easyMDE;
document.addEventListener('DOMContentLoaded', function() {
    const contentTextarea = document.getElementById('content');
    if (contentTextarea && typeof EasyMDE !== 'undefined') {
        easyMDE = new EasyMDE({
            element: contentTextarea,
            spellChecker: false,
            renderingConfig: {
                singleLineBreaks: false,
                codeSyntaxHighlighting: true,
            },
            toolbar: [
                'bold', 'italic', 'heading', '|',
                'code', 'quote', 'unordered-list', 'ordered-list', '|',
                'link', 'image', '|',
                'preview', 'side-by-side', 'fullscreen', '|',
                'guide'
            ],
            status: ['autosave', 'lines', 'words', 'cursor'],
            autofocus: false,
            placeholder: 'Write your blog post content here...',
            minHeight: '400px',
            uploadImage: true,
            imageUploadFunction: function(file, onSuccess, onError) {
                const formData = new FormData();
                formData.append('image', file);

                fetch('/api/upload/image', {
                    method: 'POST',
                    body: formData
                })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Upload failed');
                    }
                    return response.json();
                })
                .then(data => {
                    onSuccess(data.data.filePath);
                })
                .catch(error => {
                    console.error('Upload error:', error);
                    onError('Upload failed. Please try again.');
                });
            }
        });
    } else {
        console.warn('EasyMDE not loaded, using plain textarea');
    }

    // Check if editing after editor is initialized
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
                // Set content in EasyMDE if available, otherwise textarea
                if (easyMDE) {
                    easyMDE.value(post.content || '');
                } else {
                    document.getElementById('content').value = post.content || '';
                }
                document.getElementById('published').checked = post.published || false;
                document.getElementById('slug').dataset.original = post.slug || '';

                // Set publish date
                if (post.created_at) {
                    const publishDate = new Date(post.created_at).toLocaleString();
                    document.getElementById('publishDate').value = publishDate;
                } else {
                    document.getElementById('publishDate').value = 'Not published yet';
                }
            })
            .catch(error => {
                console.error('Error loading post:', error);
                alert('Error loading post data.');
            });
    }
});